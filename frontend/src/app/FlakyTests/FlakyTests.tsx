import React, { useState } from 'react';
import { Header } from '@app/utils/Header';
import { TextContent, Text, TextVariants, CardFooter, Spinner, EmptyState, EmptyStateVariant, EmptyStateIcon, Breadcrumb, BreadcrumbItem, ToolbarItem } from '@patternfly/react-core';
import { PageSection, PageSectionVariants, Title, TitleSizes } from '@patternfly/react-core';
import { DropdownItem } from '@patternfly/react-core';
import { Grid, GridItem, } from '@patternfly/react-core';
import { Card, CardTitle } from '@patternfly/react-core';
import { Toolbar, ToolbarContent } from '@patternfly/react-core';
import { getRepositoriesWithJobs, getFlakyData, getGlobalImpactData } from '../utils/APIService'
import { useSelector } from 'react-redux';
import { ExclamationCircleIcon, InfoCircleIcon } from "@patternfly/react-icons";
import { Popover } from '@patternfly/react-core';
import { Modal, Button, Flex, FlexItem } from '@patternfly/react-core';
import { Flakey, FlakyObject, TestCase } from './Types';
import { ImpactChart } from './Charts';
import { ComposableTableNestedExpandable, InnerNestedComposableTableNestedExpandable } from './Tables';
import { DropdownBasic, SpinnerBasic } from './utils';
import { useHistory } from 'react-router-dom';
import { formatDate, getRangeDates } from '@app/Reports/utils';
import { DateTimeRangePicker } from '@app/utils/DateTimeRangePicker';


const PREFERENCES_CACHE_NAME = "flakyCache"

const FlakyTests: React.FunctionComponent = () => {
  const [toggles, setToggles] = React.useState<any>([])
  const [selectedSuite, setSelectedSuite] = React.useState<string>('')
  const [data, setData] = React.useState<any>([])
  const [top10Data, setTop10Data] = React.useState<TestCase[]>([])
  const [barData, setBarData] = React.useState<any>([])
  const [selectedRepo, setSelectedRepo] = React.useState<string>('')
  const [repositories, setRepositories] = React.useState<Array<any>>([])
  const [jobs, setJobs] = React.useState<Array<any>>([])
  const [globalImpact, setGlobalImpact] = React.useState<any>([])
  const [selectedJob, setSelectedJob] = React.useState<string>('')
  const [loadingSpinner, setLoadingSpinner] = React.useState<boolean>(false)
  const currentTeam = useSelector((state: any) => state.teams.Team);
  const params = new URLSearchParams(window.location.search);
  const history = useHistory();
  const [isModalOpen, setIsModalOpen] = React.useState(false);
  const [loadingState, setLoadingState] = React.useState(false);
  const [isEmpty, setIsEmpty] = React.useState(false);
  const [rangeDateTime, setRangeDateTime] = useState(getRangeDates(15));

  // Reset rangeDateTime
  const clearRangeDateTime = () => {
    setRangeDateTime(getRangeDates(15))
  }

  const handleModalToggle = () => {
    setIsModalOpen(!isModalOpen);
  };

  const countSuiteFailures = (suites) => {
    return suites.map((suite) => {
      const c = suite.test_cases.reduce(function (acc, obj) { return acc + obj.count; }, 0);
      return { suite_name: suite.suite_name, count: c }
    })
  }

  const saveCachePreferences = () => {
    let cache = {}
    const stored = localStorage.getItem(PREFERENCES_CACHE_NAME)
    if (stored) {
      cache = JSON.parse(stored)
    }
    cache[currentTeam] = { job: selectedJob, repository: selectedRepo }
    localStorage.setItem(PREFERENCES_CACHE_NAME, JSON.stringify(cache))
  }

  const getCachePreferences = () => {
    const stored = localStorage.getItem(PREFERENCES_CACHE_NAME)
    let cache = {}
    if (stored) {
      cache = JSON.parse(stored)
    }
    return cache
  }

  React.useEffect(() => {
    setLoadingState(true)

    if (!currentTeam) { return }
    clearAllData()

    const team = params.get("team")
    const repo = params.get("repository")
    const job_name = params.get("job")
    const start = params.get("start")
    const end = params.get("end")


    getRepositoriesWithJobs(currentTeam).then(res => {
      if (res.code == 200) {
        if (res.data.length < 1 && (team == currentTeam || team == null)) {
          setLoadingState(false)
          setIsEmpty(true)
          history.push('/home/flaky?team=' + currentTeam);
        }

        setRepositories(res.data)
        const cache = getCachePreferences()

        if (repo == null || job_name == null || start == null || end == null) { // first click on page or team
          if (cache && cache[currentTeam] && cache[currentTeam].job && cache[currentTeam].repository) {
            const job = cache[currentTeam].job
            setSelectedJob(job)

            const repository = cache[currentTeam].repository
            setSelectedRepo(repository)

            const start_date = formatDate(rangeDateTime[0])
            const end_date = formatDate(rangeDateTime[1])

            history.push('/home/flaky?team=' + currentTeam + '&repository=' + repository
              + '&job=' + job + '&start=' + start_date + '&end=' + end_date)
          } else {
            console.error("Cannot get repositories and jobs")
          }
        } else {
          setSelectedRepo(repo)
          setSelectedJob(job_name)
          setRangeDateTime([new Date(start), new Date(end)])

          history.push('/home/flaky?team=' + currentTeam + '&repository=' + repo +
            '&job=' + job_name + '&start=' + start + '&end=' + end)
        }
        setLoadingState(false)
      }
    })
  }, [currentTeam]);

  React.useEffect(() => {
    if (data) {
      const organizedData = countSuiteFailures(data)
      setToggles(organizedData)
      const allCases: TestCase[] = []
      data.forEach(element => {
        allCases.push(...element.test_cases)
      })
      setTop10Data(allCases.sort((a, b) => b.count - a.count))
    }

  }, [data]);

  React.useEffect(() => {
    fillJobs(selectedRepo)
  }, [selectedRepo]);


  React.useEffect(() => {
    if (selectedJob && selectedRepo && rangeDateTime) {
      fetchData()
      saveCachePreferences()
    }

  }, [selectedJob, selectedRepo, rangeDateTime]);

  const onSuiteSelect = (value) => {
    setSelectedSuite(value)
  }

  const onDataFilter = (suite: Flakey) => {
    return suite.suite_name == selectedSuite || selectedSuite == '' || selectedSuite == 'All failures'
  }

  const clearAllData = () => {
    setBarData([])
    setSelectedJob('')
    setSelectedRepo('')
    setData([])
    setSelectedSuite('')
    setIsEmpty(false)
    clearRangeDateTime()
  }

  const fillJobs = (value) => {
    const j = repositories.filter(repo => repo.Repository.Name == value)
    if (!j || !j[0] || j[0].jobs == null) {
      setJobs([])
      return
    }
    if (j && j[0] && j[0].jobs) {
      setJobs(j[0].jobs)
    }
  }

  const onRepoSelect = (value) => {
    if (selectedRepo != value) {
      clearAllData()
      setSelectedRepo(value)
      params.set('repository', value);
      params.set('job', "");
      history.push(window.location.pathname + '?' + params.toString());
    }
  }

  const onJobSelect = (value) => {
    if (value != selectedJob) {
      setSelectedJob(value)
      params.set('job', value);
      history.push(window.location.pathname + '?' + params.toString());
    }
  }

  const fetchData = () => {
    setLoadingSpinner(true)
    if (startDate && endDate && selectedRepo && selectedRepo) {
      getFlakyData(currentTeam, selectedJob, selectedRepo, rangeDateTime, "redhat-appstudio").then(res => {
        if (res.code == 200) {
          const impact: FlakyObject = res.data
          if (impact && impact.suites) {
            setData(res.data.suites)
          }
          if (impact && impact.global_impact) {
            const gd = [{ Date: rangeDateTime[0].toISOString().split('T')[0], global_impact: res.data.global_impact }, { Date: rangeDateTime[1].toISOString().split('T')[0], global_impact: res.data.global_impact }]
            setGlobalImpact(gd)
          }
          setLoadingSpinner(false)
        } else {
          console.log("error", res)
        }
      })
      getGlobalImpactData(currentTeam, selectedJob, selectedRepo, rangeDateTime, "redhat-appstudio").then(res => {
        if (res.code == 200) {
          const impact: any = res.data
          if (impact && impact.length > 0) {
            setBarData(impact.map(impact => { impact.Date = impact.date.split(' ')[0].split('T')[0]; return impact; }))
          }
          setLoadingSpinner(false)
        } else {
          console.log("error", res)
        }
      })
    }
  }

  function handleChange(event, from, to) {
    setRangeDateTime([from, to])
    params.set("start", formatDate(from))
    params.set("end", formatDate(to))
    history.push(window.location.pathname + '?' + params.toString());
  }

  const startDate = rangeDateTime[0]
  const endDate = rangeDateTime[1]


  return (
    <React.Fragment>
      <Header info="Observe the impact of the flaky tests that are affecting CI."></Header>
      <div style={{ marginTop: 15, marginBottom: 15, marginLeft: 15 }}>
        <Breadcrumb>
          <BreadcrumbItem>OpenShift CI</BreadcrumbItem>
          <BreadcrumbItem>Flaky Tests</BreadcrumbItem>
        </Breadcrumb>
      </div>
      <PageSection variant={PageSectionVariants.light}>
        <Title headingLevel="h3" size={TitleSizes['2xl']}>
          Flaky tests impacting CI
        </Title>
      </PageSection>
      <PageSection style={{ minHeight: '120vh' }} variant={PageSectionVariants.default}>
        <div>
          {loadingState && <div style={{ width: '100%', textAlign: "center" }}>
            <Spinner isSVG diameter="80px" aria-label="Contents of the custom size example" style={{ margin: "100px auto" }} />
          </div>
          }
          {isEmpty && !loadingState && <EmptyState variant={EmptyStateVariant.xl}>
            <EmptyStateIcon icon={ExclamationCircleIcon} />
            <Title headingLevel="h1" size="lg">
              No data found for this team.
            </Title>
          </EmptyState>
          }

          {!loadingState && !isEmpty &&
            (
              <Grid hasGutter>

                <GridItem span={12}>
                  <TextContent className='bg-white'>
                    <Text component={TextVariants.h1}>Flaky tests</Text>
                  </TextContent>
                </GridItem>
                <GridItem span={12}>
                  <Toolbar id="toolbar-items">
                    <ToolbarContent>
                      <ToolbarItem>
                        <DropdownBasic selected={selectedRepo} toggles={
                          repositories.map((repo, idx) => <DropdownItem key={idx} name={repo.Repository.Name}> {repo.Repository.Name} </DropdownItem>)
                        } onSelect={onRepoSelect} placeholder="Select a repository"></DropdownBasic>
                      </ToolbarItem>
                      <ToolbarItem>
                        <DropdownBasic selected={selectedJob} toggles={
                          jobs.map((job, idx) => <DropdownItem key={idx} name={job}> {job} </DropdownItem>)
                        } onSelect={onJobSelect} placeholder="Select a job"></DropdownBasic>
                      </ToolbarItem>
                      <ToolbarItem style={{ minWidth: "fitContent", maxWidth: "fitContent" }}>
                        <DateTimeRangePicker
                          startDate={startDate}
                          endDate={endDate}
                          handleChange={(event, from, to) => handleChange(event, from, to)}
                        >
                        </DateTimeRangePicker>
                      </ToolbarItem>
                    </ToolbarContent>
                  </Toolbar>
                </GridItem>
                <GridItem>
                  <Grid hasGutter className='bg-white'>
                    <GridItem span={11}>
                      <div>
                        <Title headingLevel="h3">
                          <Popover
                            headerContent={<div>What is this chart about?</div>}
                            bodyContent={
                              <div>
                                <Title headingLevel="h2">Global impact</Title>
                                <Text>Global impact shows the average of the tests impact in the selected time range.</Text>
                              </div>
                            }
                          >
                            <InfoCircleIcon></InfoCircleIcon>
                          </Popover>
                          <span style={{ paddingLeft: '1em' }}>Impact on CI suite (%)</span>
                        </Title>
                        <SpinnerBasic isLoading={loadingSpinner}></SpinnerBasic>
                        <ImpactChart data={barData} x="Date" y="global_impact" secondaryData={globalImpact}></ImpactChart>
                      </div>
                    </GridItem>

                  </Grid>
                </GridItem>
                <GridItem style={{ clear: 'both', minHeight: '1em' }} span={12}>
                </GridItem>
                <GridItem span={12}>
                  <Toolbar id="toolbar-items" style={{ padding: '1em' }}>
                    <Flex>
                      <FlexItem>
                        <DropdownBasic selected={selectedSuite} toggles={
                          [<DropdownItem key={'all'} name="All failures"> All failures</DropdownItem>,
                          ...toggles.sort((a, b) => b.count - a.count).map((toggle, idx) => <DropdownItem key={idx} name={toggle.suite_name}> {toggle.suite_name} (<strong style={{ color: 'red' }}>{toggle.count}</strong>) </DropdownItem>)
                          ]
                        } onSelect={onSuiteSelect} placeholder="Select a failing suite"></DropdownBasic>
                      </FlexItem>
                      <FlexItem align={{ default: 'alignRight' }}>
                        <Button variant="primary" onClick={handleModalToggle} ouiaId="ShowBasicModal" style={{ margin: "auto 2em" }}>
                          Show top 20 failing test cases
                        </Button>
                      </FlexItem>
                    </Flex>
                  </Toolbar>
                  <Modal
                    title="Basic modal"
                    isOpen={isModalOpen}
                    width="80%"
                    onClose={handleModalToggle}
                    actions={[]}
                    ouiaId="BasicModal"
                  >
                    <Card className='card-no-border' style={{ padding: "1em" }}>
                      <CardTitle>Top 20 Failing test cases</CardTitle>
                      <InnerNestedComposableTableNestedExpandable test_cases={top10Data.slice(0, 20)} rowIndex={0}></InnerNestedComposableTableNestedExpandable>
                      <CardFooter></CardFooter>
                    </Card>
                  </Modal>
                </GridItem>
                <GridItem span={12}>
                  <SpinnerBasic isLoading={loadingSpinner}></SpinnerBasic>
                  <ComposableTableNestedExpandable teams={data.filter(onDataFilter)}></ComposableTableNestedExpandable>
                </GridItem>
              </Grid>
            )
          }
        </div>
      </PageSection >
    </React.Fragment >
  )
}

export { FlakyTests };
