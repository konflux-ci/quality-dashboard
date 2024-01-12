import React, { useEffect, useState } from 'react';
import { CubesIcon, ExclamationCircleIcon, InfoCircleIcon, OpenDrawerRightIcon } from '@patternfly/react-icons';
import {
  PageSection, PageSectionVariants,
  EmptyState,
  EmptyStateVariant,
  EmptyStateIcon,
  EmptyStateBody,
  Title, TitleSizes,
  Alert, AlertGroup, AlertActionCloseButton,
  Spinner, Breadcrumb, BreadcrumbItem, Tooltip,
  Drawer,
  DrawerContent,
  DrawerPanelContent,
  DrawerHead,
  DrawerActions,
  DrawerCloseButton,
  TextContent,
  TextList,
  TextListItem
} from '@patternfly/react-core';
import { Toolbar, ToolbarItem, ToolbarContent } from '@patternfly/react-core';
import { Button } from '@patternfly/react-core';
import { Select, SelectOption, SelectVariant } from '@patternfly/react-core';
import { getFlakyData, getProwJobStatistics, getTeams, getJobNamesAndTypes, listTeamRepos, getProwJobMetricsDaily } from '@app/utils/APIService';
import { Grid, GridItem } from '@patternfly/react-core';
import {
  DashboardLineChart,
  DashboardLineChartData,
  JobMetric
} from '@app/utils/sharedComponents';
import { Card, CardTitle, CardBody, CardFooter } from '@patternfly/react-core';
import { Flex, FlexItem } from '@patternfly/react-core';

import { Divider, TextVariants, Text } from '@patternfly/react-core';
import { ReactReduxContext, useSelector } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { formatDate, getRangeDates } from './utils';
import { DateTimeRangePicker } from '../utils/DateTimeRangePicker';
import { validateRepositoryParams, validateParam } from '@app/utils/utils';
import { Header } from '@app/utils/Header';

// eslint-disable-next-line prefer-const
let Reports = () => {

  const [prowVisible, setProwVisible] = useState(false);
  const [loadingState, setLoadingState] = useState(false);
  const [noData, setNoData] = useState(false);
  const [alerts, setAlerts] = React.useState<React.ReactNode[]>([]);
  const [impact, setImpact] = useState("");
  const { store } = React.useContext(ReactReduxContext);
  const state = store.getState();

  /* 
  Toolbar dropdowns logic and helpers
  */

  const [repositories, setRepositories] = useState<any[]>([]);
  const [repoName, setRepoName] = useState("");
  const [repoNameFormatted, setRepoNameFormatted] = useState("");
  const [repoOrg, setRepoOrg] = useState("");
  const [jobType, setjobType] = useState("");
  const [jobMeta, setJobMeta] = useState<any[]>([]);
  const [jobTypes, setJobTypes] = useState<string[]>([]);
  const [jobName, setjobName] = useState("");
  const [jobNames, setJobNames] = useState<string[]>([]);
  const [jobTypeToggle, setjobTypeToggle] = useState(false);
  const [jobNameToggle, setjobNameToggle] = useState(false);
  const [repoNameToggle, setRepoNameToggle] = useState(false);
  const currentTeam = useSelector((state: any) => state.teams.Team);
  const history = useHistory();
  const params = new URLSearchParams(window.location.search);
  const [rangeDateTime, setRangeDateTime] = useState(getRangeDates(10));
  const [isInvalid, setIsInvalid] = useState(false);
  const [isExpanded, setIsExpanded] = React.useState(false);
  const drawerRef = React.useRef<HTMLDivElement>();

  // Reset all dropdowns and state variables
  const clearAll = () => {
    setProwVisible(false); // hide the dashboard leaving only the toolbar
    setNoData(false)
    clearJobType()
    clearRepo()
    clearRangeDateTime()
    setIsInvalid(false)
    clearJobName()
  }

  // Reset params
  const clearParams = () => {
    clearAll()
    history.push(window.location.pathname + '?' + "team=" + params.get("team"));
  }

  // Reset the repository dropdown
  const clearRepo = () => {
    setRepoName("")
    setRepoNameFormatted("")
    setRepoOrg("")
    setRepoNameToggle(false)
  }

  // Reset the jobType dropdown
  const clearJobType = () => {
    setjobType("");
    setjobTypeToggle(false);
    setjobNameToggle(false);
  }

  const clearJobName = () => {
    setjobName("")
  }

  // Reset rangeDateTime
  const clearRangeDateTime = () => {
    setRangeDateTime(getRangeDates(10))
  } 
  
  const onExpand = () => {
    drawerRef.current && drawerRef.current.focus();
  };

  const onClick = () => {
    setIsExpanded(!isExpanded);
  };

  const onCloseClick = () => {
    setIsExpanded(false);
  };

  const panelContent = (
    <DrawerPanelContent isResizable defaultSize={'500px'} minSize={'150px'}>
      <DrawerHead>
        <>
          <TextContent>
            <Title headingLevel="h1">Test Reports</Title>
            <Text>
              This page aims help you analyze the OpenShift CI prow jobs executions in your team&apos;s GitHub repositories.
            </Text>
            <Title headingLevel="h1">Repositories</Title>
            <Text>
              Lists all the repositories where the prow jobs runs against.
            </Text>
            <Title headingLevel="h1">Jobs</Title>
            <Text>
              Lists all the jobs executed in the target GitHub Repository. Example: pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests.
            </Text>
            <Title headingLevel="h1">Job Types</Title>
            <TextList>
              <TextListItem>presubmit - executes every time a PR is updated.</TextListItem>
              <TextListItem>periodic - executes periodically.</TextListItem>
              <TextListItem>postsubmit - executes after a PR is merged.</TextListItem>
            </TextList>
          </TextContent>
        </>
        <DrawerActions>
          <DrawerCloseButton onClick={onCloseClick} />
        </DrawerActions>
      </DrawerHead>
    </DrawerPanelContent>
  );

  // Called onChange of the repository dropdown element. This set repository name and organization state variables, or clears them when placeholder is selected
  const setRepoNameOnChange = (event, selection, isPlaceholder) => {
    console.log(selection, repositories[selection])
    if (isPlaceholder) {
      clearRepo()
    }
    else {
      setRepoName(repositories[selection].Repository.Name);
      setRepoNameFormatted(repositories[selection].Repository.Name);
      setRepoOrg(repositories[selection].Repository.Owner.Login);
      setRepoNameToggle(false)
      params.set("repository", repositories[selection].Repository.Name)
      params.set("organization", repositories[selection].Repository.Owner.Login)

      getJobNamesAndTypes(repositories[selection].Repository.Name, repositories[selection].Repository.Owner.Login)
        .then((data: any) => {
          setJobMeta(data)
          setJobTypes(data.map(el => el.job_type).filter((value, index, self) => self.indexOf(value) === index))
          setJobNames(data.map(el => el.job_name).filter((value, index, self) => self.indexOf(value) === index))
          setjobType(data[0].job_type)
          setjobName(data[0].job_name)
          params.set("job_type", data[0].job_type)
          params.set("job_name", data[0].job_name)
          params.set("start", formatDate(rangeDateTime[0]))
          params.set("end", formatDate(rangeDateTime[1]))
          history.push(window.location.pathname + '?' + params.toString());
        }
      );
    }
  };

  // Called onChange of the jobType dropdown element. This set repository name and organization state variables, or clears them when placeholder is selected

  const setjobTypeOnChange = (event, selection, isPlaceholder) => {
    if (isPlaceholder) {
      clearJobType()
    }
    else {
      clearJobName()
      setjobType(selection);
      setjobTypeToggle(false);
      setjobNameToggle(false);
      setIsInvalid(false)
      params.set("job_type", selection)
      params.delete("job_name")
      history.push(window.location.pathname + '?' + params.toString());
    }
  };

  // Called onChange of the jobName dropdown element.
  const setjobNameOnChange = (event, selection, isPlaceholder) => {
    if (isPlaceholder) {
      clearJobName()
    }
    else {
      setjobName(selection);
      setjobTypeToggle(false);
      setjobNameToggle(false);
      setIsInvalid(false)
      params.set("job_name", selection)
      history.push(window.location.pathname + '?' + params.toString());
    }
  };

  // Filter for the job names based on the jobType currently selected
  const filterJobNames = (value) => {
    if(!jobType || jobType == "" || jobType == null){
      return true
    } 
    return jobMeta.filter(j => j.job_type == jobType).map(el => el.job_name).includes(value)
  }

  /* 
  ProwJobs logic to populate dashboard
  */

  const [prowJobsStats, setprowJobsStats] = useState<any | null>(null);
  const [prowJobMetrics, setprowJobMetrics] = useState<DashboardLineChartData | null>(null);
  const [prowJobFailuerMetrics, setprowJobFailureMetrics] = useState<DashboardLineChartData | null>(null);
  //const [prowJobs, setProwJobs] = useState<Job[] | null>(null);

  // Get the prow jobs from API
  const getProwJob = async () => {
    // Hide components and show loading spinner 
    setProwVisible(false)
    setLoadingState(true)
    setNoData(false)
    try {
      // Get statistics and metrics
      const stats = await getProwJobStatistics(repoName, repoOrg, jobType, jobName, rangeDateTime)
      setprowJobsStats(stats)

      // Get jobs
      //const prowJobs = await getProwJobs(repoName, repoOrg, rangeDateTime)
      //setProwJobs(prowJobs)

      // Set UI for showing data and disable spinner
      setLoadingState(false)
      setProwVisible(true)
    }
    catch (e) {
      // Set UI to empty page
      setProwVisible(false);
      setLoadingState(false);

      // Show error alert
      if (e != "No jobs detected in OpenShift CI") {
        setAlerts(prevAlerts => {
          return [...prevAlerts,
          <Alert
            variant="danger"
            timeout={5000}
            title="Error fetching data from server"
            key={0}
            actionClose={
              <AlertActionCloseButton
                title="Error fetching data"
                variantLabel={`danger alert`}
                onClose={() => setAlerts([])}
              />
            }
          />
          ]
        });
      } else {
        // Set UI to no data page
        setNoData(true)
      }
    }
  }

  // Prepare data for the line chart
  

  function handleChange(event, from, to) {
    setRangeDateTime([from, to])
    params.set("start", formatDate(from))
    params.set("end", formatDate(to))
    history.push(window.location.pathname + '?' + params.toString());
  }

  const start = rangeDateTime[0]
  const end = rangeDateTime[1]

  const getChartData = async () => {
    let data = await getProwJobMetricsDaily(repoName, repoOrg, jobType, jobName, rangeDateTime)
    // sort the metrics by date
    data = data.sort(function(a,b){
      // Turn your strings into dates, and then subtract them
      // to get a value that is either negative, positive, or zero.
      return new Date(a.start_date).valueOf() - new Date(b.start_date).valueOf();
    });

    console.log(data)
    const beautifiedData: DashboardLineChartData = {
      "FAILED_JOB_RUNS": { data: [], style: { data: { stroke: "rgba(255, 0, 0, 0.6)", strokeWidth: 2 } } },
      "SUCCESS_JOB_RUNS": { data: [], style: { data: { stroke: "rgba(60, 179, 113, 0.6)", strokeWidth: 2 } } },
    };
    
    const beautifiedFailureData: DashboardLineChartData = {
      "EXTERNAL_SERVICES_IMPACT": { data: [], style: { data: { stroke: "rgba(255, 165, 0, 0.7)", strokeWidth: 2 } } },
      "FLAKY_TESTS_IMPACT": { data: [], style: { data: { stroke: "rgba(255, 0, 0, 0.7)", strokeWidth: 2 } } },
      "INFRASTRUCTURE_IMPACT": { data: [], style: { data: { stroke: "rgba(106, 90, 205, 0.7)", strokeWidth: 2 } } },
      "UNKNOWN_FAILURES_IMPACT": { data: [], style: { data: { stroke: "rgba(180, 180, 180, 0.7)", strokeWidth: 2 } } },
    };

    data.map(jobMetric => {
      beautifiedData["FAILED_JOB_RUNS"].data.push({ name: 'failed_job_runs_%', x: new Date(jobMetric.start_date).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: jobMetric.jobs_runs.failed_percentage })
      beautifiedData["SUCCESS_JOB_RUNS"].data.push({ name: 'success_job_runs_%', x: new Date(jobMetric.start_date).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: jobMetric.jobs_runs.success_percentage })

      beautifiedFailureData["EXTERNAL_SERVICES_IMPACT"].data.push({ name: 'external_services_impact_%', x: new Date(jobMetric.start_date).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: jobMetric.jobs_impacts.external_services_impact.percentage })
      beautifiedFailureData["FLAKY_TESTS_IMPACT"].data.push({ name: 'flaky_tests_impact_%', x: new Date(jobMetric.start_date).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: jobMetric.jobs_impacts.flaky_tests_impact.percentage })
      beautifiedFailureData["INFRASTRUCTURE_IMPACT"].data.push({ name: 'infrasctructure_impact_%', x: new Date(jobMetric.start_date).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: jobMetric.jobs_impacts.infrastructure_impact.percentage })
      beautifiedFailureData["UNKNOWN_FAILURES_IMPACT"].data.push({ name: 'unknow_failures_impact_%', x: new Date(jobMetric.start_date).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: jobMetric.jobs_impacts.unknown_failures_impact.percentage })

    })

    setprowJobMetrics(beautifiedData)
    setprowJobFailureMetrics(beautifiedFailureData)

  }

  // when a job name is selected, get the job data from api
  useEffect(() => {
    if(currentTeam && jobName && repoName && start.toISOString() && end.toISOString() && repoOrg){
      getProwJob();
      getFlakyData(currentTeam, jobName, repoName, start.toISOString(), end.toISOString(), repoOrg).then(res => { setImpact(res.data.global_impact)});
      getChartData()
    } else {
      setImpact("")
    }
  }, [jobName]);

  useEffect(() => {
    clearJobName()
    clearJobType()
  }, [repoName]);

  // When component is mounted, get the list of repo and orgs from API and populate the dropdowns
  useEffect(() => {
    setLoadingState(true)
    const team = params.get("team")

    if ((team != null) && (team != state.teams.Team)) {
      getTeams().then(res => {
        if (!validateParam(res.data, team)) {
          setLoadingState(false)
          setIsInvalid(true)
        }
      })
    }

    if (state.teams.Team != "") {
      setRepositories([])
      clearAll()

      const repository = params.get("repository")
      const organization = params.get("organization")
      const job_type = params.get("job_type")
      const job_name = params.get("job_name")
      const start = params.get("start")
      const end = params.get("end")

      listTeamRepos(state.teams.Team)
        .then((data: any) => {
          let dropDescr = ""
          if (data.length < 1 && (team == state.teams.Team || team == null)) {
            dropDescr = "No Repositories"
            setLoadingState(false)
            history.push('/reports/test?team=' + currentTeam)
          }
          else { dropDescr = "Select a repository" }

          if (data.length > 0 && (team == state.teams.Team || team == null)) {
            setRepositories(data)

            if (repository && organization && job_type && job_name && start && end) {

              if (validateRepositoryParams(data, repository, organization)) {
                setRepoName(repository)
                setRepoNameFormatted(repository)
                setRepoOrg(organization)
                setRangeDateTime([new Date(start), new Date(end)])
                getJobNamesAndTypes(repository, organization)
                  .then((data: any) => {
                    setJobMeta(data)
                    setJobTypes(data.map(el => el.job_type).filter((value, index, self) => self.indexOf(value) === index))
                    setJobNames(data.map(el => el.job_name).filter((value, index, self) => self.indexOf(value) === index))
                    setjobName(job_name)
                    setjobType(job_type)
                  }
                );
              } else {
                setLoadingState(false)
                setIsInvalid(true)
              }

            } 
            else {
              setRepoName(data[0].Repository.Name)
              setRepoNameFormatted(data[0].Repository.Name)
              setRepoOrg(data[0].Repository.Owner.Login)
              setjobType("presubmit") // all repos in OpenShift CI have presubmit type job

              getJobNamesAndTypes(data[0].Repository.Name, data[0].Repository.Owner.Login)
                .then((data: any) => {
                  setJobMeta(data)
                  setJobTypes(data.map(el => el.job_type).filter((value, index, self) => self.indexOf(value) === index))
                  setJobNames(data.map(el => el.job_name).filter((value, index, self) => self.indexOf(value) === index))
                  setjobName(data[0].job_name)
                  setjobType(data[0].job_type)
                }
              );

              const start_date = formatDate(rangeDateTime[0])
              const end_date = formatDate(rangeDateTime[1])

              history.push('/reports/test?team=' + currentTeam + '&organization=' + data[1].organization + '&repository=' + data[1].repoName
                + '&job_type=presubmit' + '&start=' + start_date + '&end=' + end_date)
              setLoadingState(false)
            }

          }
        })
    }
  }, [setRepositories, currentTeam]);

  return (

    <React.Fragment>
      {/* page title bar */}
      <Header info="Observe the CI metrics of all your CI jobs."></Header>
      <div style={{ padding: 15, background: "white", border: "1px solid lightgrey" }}>
        <Breadcrumb>
          <BreadcrumbItem>OpenShift CI</BreadcrumbItem>
          <BreadcrumbItem>Tests Reports</BreadcrumbItem>
        </Breadcrumb>
      </div>
      <PageSection variant={PageSectionVariants.light}>
        <Title headingLevel="h3" size={TitleSizes['2xl']}>
          Tests Reports
          <Button onClick={onClick} variant="link" icon={<OpenDrawerRightIcon />} iconPosition="right">
            Learn more
          </Button>
        </Title>
      </PageSection>
      <Drawer isExpanded={isExpanded} onExpand={onExpand}>
        <DrawerContent panelContent={panelContent}>
          <PageSection variant={PageSectionVariants.light} style={{ padding: "10px 5px", background: "white" }}>
            {/* the following toolbar will contain the form (dropdowns and button) to request data to the server */}
            <Toolbar style={{ width: prowVisible ? '100%' : '100%', margin: prowVisible ? 'auto' : '0', padding: 0 }}>
              <ToolbarContent style={{ textAlign: 'left' }}>
                <ToolbarItem style={{ minWidth: "fitContent", maxWidth: "fitContent" }}>
                  <Select
                    width={600}
                    variant={SelectVariant.typeahead}
                    typeAheadAriaLabel="Select a repository"
                    isOpen={repoNameToggle}
                    onToggle={setRepoNameToggle}
                    selections={repoNameFormatted}
                    onSelect={setRepoNameOnChange}
                    onClear={clearAll}
                    aria-labelledby="typeahead-select"
                    placeholderText="Select a repository"
                  >
                    {repositories.map((value, index) => (
                      <SelectOption key={index} value={index} description={value.Repository.Name + "/" + value.Repository.Owner.Login} isDisabled={value.isPlaceholder}>{value.Repository.Name}</SelectOption>
                    ))}
                  </Select>
                </ToolbarItem>
                <ToolbarItem style={{ minWidth: "fitContent", maxWidth: "fitContent", border: 'none' }}>
                  <DateTimeRangePicker
                    startDate={start}
                    endDate={end}
                    handleChange={(event, from, to) => handleChange(event, from, to)}
                  >
                  </DateTimeRangePicker>
                </ToolbarItem>
              </ToolbarContent>
              <ToolbarContent style={{ textAlign: 'left' }}>
                {
                  <ToolbarItem>
                    <Select placeholderText="Filter by job type" isOpen={jobTypeToggle} onToggle={setjobTypeToggle} selections={jobType} onSelect={setjobTypeOnChange} aria-label="Select Input">
                      {jobTypes.map((value, index) => (
                        <SelectOption key={index} value={value}>{value}</SelectOption>
                      ))}
                    </Select>
                  </ToolbarItem>
                }
                {
                  <ToolbarItem>
                    <Select placeholderText="Filter by job name" width={600} isOpen={jobNameToggle} onToggle={setjobNameToggle} selections={jobName} onSelect={setjobNameOnChange} aria-label="Select Input">
                      {jobNames.filter(filterJobNames).map((value, index) => (
                        <SelectOption key={index} isPlaceholder={false} value={value}>{value}</SelectOption>
                      ))}
                    </Select>
                  </ToolbarItem>
                }
                <ToolbarItem style={{ minWidth: "fitContent", maxWidth: "fitContent" }}>
                  <Button variant="link" onClick={clearParams}>Clear all selections</Button>
                </ToolbarItem>
              </ToolbarContent>
            </Toolbar>
          </PageSection>
          {/* main content  */}
          <PageSection>
            {/* alertGroup will show toast notification (on the top left) when an error occurs */}
            <AlertGroup isToast isLiveRegion> {alerts} </AlertGroup>
            {/* if the server has not provided any data or if the clear button is clicked or if the page is in its initial state, this empty placeholder will be shown */}
            {loadingState && <div style={{ width: '100%', textAlign: "center" }}>
              <Spinner isSVG diameter="80px" aria-label="Contents of the custom size example" style={{ margin: "100px auto" }} />
            </div>
            }
            {window.location.pathname == '/reports/test' && params.get("team") != null && params.get("organization") == null && !loadingState && <EmptyState variant={EmptyStateVariant.xl}>
              <EmptyStateIcon icon={CubesIcon} />
              <Title headingLevel="h1" size="lg">
                No job selected yet.
              </Title>
              <EmptyStateBody>
                Please select a repository and an organization to see the last job&apos;s details
              </EmptyStateBody>
            </EmptyState>
            }
            {!isInvalid && noData && !loadingState && <EmptyState variant={EmptyStateVariant.xl}>
              <EmptyStateIcon icon={ExclamationCircleIcon} />
              <Title headingLevel="h1" size="lg">
                No jobs detected in OpenShift CI.
              </Title>
            </EmptyState>
            }
            {isInvalid && !loadingState && <EmptyState variant={EmptyStateVariant.xl}>
              <EmptyStateIcon icon={ExclamationCircleIcon} />
              <Title headingLevel="h1" size="lg">
                Something went wrong. Please, check the URL.
              </Title>
            </EmptyState>
            }
            {/* this section will show statistics and details about job and suites */}
            <React.Fragment>
              {prowVisible && <div>
                {/* this section will show the job's chart over time and last execution stats */}
                {prowJobsStats !== null &&
                  <Grid hasGutter>
                    <GridItem span={5} rowSpan={1}>
                      <Card style={{minHeight: "30vh"}}>
                        <CardTitle>Jobs Executed</CardTitle>
                        <CardBody>
                          <Text component={TextVariants.p} style={{minHeight: "2em"}}>Number of jobs executed in the selected time range, with success and failure rate. </Text>
                          <Flex className="example-border" justifyContent={{ default: 'justifyContentSpaceEvenly' }} flexWrap={{ default: 'nowrap' }} direction={{ default: 'column', sm: 'row' }}>
                            <FlexItem>
                              <Card style={{ border: 'none', boxShadow: "none", textAlign: "center" }}>
                                <CardTitle>
                                  <Title headingLevel="h1" size={TitleSizes['2xl']}>
                                    {prowJobsStats.jobs_runs ? prowJobsStats.jobs_runs.total : "-"}
                                  </Title>
                                </CardTitle>
                                <CardBody>
                                  <Title headingLevel="h1" style={{minHeight: "2em"}}>Total</Title>
                                </CardBody>
                                <CardFooter style={{ color: "black" }}>
                                  <div style={{minHeight: "3em"}}>Job Runs</div>
                                  <Tooltip content={<div>Total number of jobs executed in the selected time range.</div>}>
                                    <InfoCircleIcon></InfoCircleIcon>
                                  </Tooltip>
                                </CardFooter>
                              </Card>
                            </FlexItem>
                            <Divider orientation={{ default: 'vertical' }} />
                            <FlexItem>
                              <Card style={{ border: 'none', boxShadow: "none", textAlign: "center", color: "#1E4F18" }}>
                                <CardTitle>
                                  <Title headingLevel="h1" size={TitleSizes['2xl']}>
                                    {prowJobsStats.jobs_runs ? prowJobsStats.jobs_runs.success : "-"}
                                  </Title>
                                </CardTitle>
                                <CardBody>
                                  <Title headingLevel="h1" style={{minHeight: "2em"}}>
                                    {prowJobsStats.jobs_runs ? prowJobsStats.jobs_runs.success_percentage : "-"}%
                                  </Title>
                                </CardBody>
                                <CardFooter style={{ color: "black" }}>
                                  <div style={{minHeight: "3em"}}>Completed Jobs</div>
                                  <Tooltip content={<div>Count and percentage of jobs that completed successfully.</div>}>
                                    <InfoCircleIcon></InfoCircleIcon>
                                  </Tooltip>
                                </CardFooter>
                              </Card>
                            </FlexItem>
                            <Divider orientation={{ default: 'vertical' }} />
                            <FlexItem>
                              <Card style={{ border: 'none', boxShadow: "none", textAlign: "center", color: "#A30000" }}>
                                <CardTitle>
                                  <Title headingLevel="h1" size={TitleSizes['2xl']}>
                                    {prowJobsStats.jobs_runs ? prowJobsStats.jobs_runs.failures : "-"}
                                  </Title>
                                </CardTitle>
                                <CardBody>
                                  <Title headingLevel="h1" style={{minHeight: "2em"}}>
                                    {prowJobsStats.jobs_runs ? prowJobsStats.jobs_runs.failed_percentage : "-"}%
                                  </Title>
                                </CardBody>
                                <CardFooter style={{ color: "black" }}>
                                  <div style={{minHeight: "3em"}}>Failed Jobs</div>
                                  <Tooltip content={<div>Count and percentage of jobs that failed. See the Failures card to understand the most relevant reasons of failure.</div>}>
                                    <InfoCircleIcon></InfoCircleIcon>
                                  </Tooltip>
                                </CardFooter>
                              </Card>
                            </FlexItem>
                            <FlexItem></FlexItem>
                          </Flex>
                        </CardBody>
                      </Card>
                    </GridItem>
                    <GridItem span={7} rowSpan={1}>
                      <Card style={{minHeight: "30vh"}}>
                        <CardTitle>Failures</CardTitle>
                        <CardBody>
                          <Text component={TextVariants.p} style={{minHeight: "2em"}}>The percentage of failures grouped by most common reasons, considering the failed jobs in the selected time range.</Text>
                          <Flex className="example-border" justifyContent={{ default: 'justifyContentSpaceBetween' }} flexWrap={{ default: 'nowrap' }} direction={{ default: 'column', sm: 'row' }}>
                            <FlexItem></FlexItem>
                            <FlexItem>
                              <Card style={{ border: 'none', boxShadow: "none", textAlign: "center", color: "#A30000" }}>
                                <CardTitle>
                                  <Title headingLevel="h1" size={TitleSizes['2xl']}>
                                    {prowJobsStats.jobs_impacts ? prowJobsStats.jobs_impacts.infrastructure_impact.total : "-"}
                                  </Title>
                                </CardTitle>
                                <CardBody>
                                  <Title headingLevel="h1" style={{minHeight: "2em"}}>
                                    {prowJobsStats.jobs_impacts ? prowJobsStats.jobs_impacts.infrastructure_impact.percentage : "-"}%
                                  </Title>
                                </CardBody>
                                <CardFooter style={{ color: "black" }}>
                                  <div style={{minHeight: "3em"}}>CI Fail</div>
                                  <Tooltip content={<div>Count and percentage of jobs that failed due to CI infrastructure (cases where the job was not scheduled, for example), considering the failing jobs in the selected period of time.</div>}>
                                    <InfoCircleIcon></InfoCircleIcon>
                                  </Tooltip>
                                </CardFooter>
                              </Card>
                            </FlexItem>
                            <Divider orientation={{ default: 'vertical' }} />
                            <FlexItem>
                              <Card style={{ border: 'none', boxShadow: "none", textAlign: "center", color: "#A30000" }}>
                                <CardTitle>
                                  <Title headingLevel="h1" size={TitleSizes['2xl']}>
                                  {prowJobsStats.jobs_impacts ? prowJobsStats.jobs_impacts.external_services_impact.total : "-"}
                                  </Title>
                                </CardTitle>
                                <CardBody>
                                  <Title headingLevel="h1" style={{minHeight: "2em"}}> 
                                  {prowJobsStats.jobs_impacts ? prowJobsStats.jobs_impacts.external_services_impact.percentage : "-"}%
                                  </Title>
                                </CardBody>
                                <CardFooter style={{ color: "black" }}>
                                  <div style={{minHeight: "3em"}}>External Services Outage</div>
                                  <Tooltip content={<div>Count and percentage of jobs that failed due to external services outage (like Github, Quay.io, etc.), considering the failing jobs in the selected period of time.</div>}>
                                    <InfoCircleIcon></InfoCircleIcon>
                                  </Tooltip>
                                </CardFooter>
                              </Card>
                            </FlexItem>
                            <Divider orientation={{ default: 'vertical' }} />
                            <FlexItem>
                              <Card style={{ border: 'none', boxShadow: "none", textAlign: "center", color: "#A30000" }}>
                                <CardTitle>
                                  <Title headingLevel="h1" size={TitleSizes['2xl']}>
                                    {prowJobsStats.jobs_impacts ? prowJobsStats.jobs_impacts.flaky_tests_impact.total : "-"}
                                  </Title>
                                </CardTitle>
                                <CardBody>
                                  <Title headingLevel="h1" style={{minHeight: "2em"}}>
                                    {prowJobsStats.jobs_impacts ? prowJobsStats.jobs_impacts.flaky_tests_impact.percentage : "-"}%
                                  </Title>
                                </CardBody>
                                <CardFooter style={{ color: "black" }}>
                                  <div style={{minHeight: "3em"}}>
                                    {prowJobsStats ? <a href={'/home/flaky?team=' + currentTeam + '&repository=' + repoName + '&job=' + jobName + '&start=' + start.toISOString() + '&end=' + end.toISOString()}>Flaky Tests</a> : "Flaky Tests"}
                                  </div>
                                  <Tooltip content={<div>Count and percentage of jobs that failed due to flaky tests, considering the failing jobs in the selected period of time. See the flaky tests page to see more details.</div>}>
                                    <InfoCircleIcon></InfoCircleIcon>
                                  </Tooltip>
                                </CardFooter>
                              </Card>
                            </FlexItem>
                            <Divider orientation={{ default: 'vertical' }} />
                            <FlexItem>
                              <Card style={{ border: 'none', boxShadow: "none", textAlign: "center", color: "#A30000" }}>
                                <CardTitle >
                                  <Title headingLevel="h1" size={TitleSizes['2xl']}>
                                    {prowJobsStats.jobs_impacts ? prowJobsStats.jobs_impacts.unknown_failures_impact.total : "-"}
                                  </Title>
                                </CardTitle>
                                <CardBody>
                                  <Title headingLevel="h1" style={{minHeight: "2em"}}>
                                    {prowJobsStats.jobs_impacts ? prowJobsStats.jobs_impacts.unknown_failures_impact.percentage : "-"}%
                                  </Title>
                                </CardBody>
                                <CardFooter style={{ color: "black" }}>
                                  <div style={{minHeight: "3em"}}>Other reasons</div>
                                  <Tooltip content={<div>Count and percentage of jobs that failed due to other reasons (like RHTAP installation failures, cluster build errors, etc.) considering the failing jobs in the selected period of time.</div>}>
                                    <InfoCircleIcon></InfoCircleIcon>
                                  </Tooltip>
                                </CardFooter>
                              </Card>
                            </FlexItem>
                            <FlexItem></FlexItem>
                          </Flex>
                        </CardBody>
                      </Card>
                    </GridItem>
                  </Grid>
                }
                {
                prowJobMetrics !== null && <Grid hasGutter style={{ margin: "20px 0px" }} sm={6} md={4} lg={3} xl2={1}>
                  <GridItem span={12} rowSpan={5}><DashboardLineChart data={prowJobMetrics}></DashboardLineChart></GridItem>
                </Grid>
                }
                {
                prowJobFailuerMetrics !== null && <Grid hasGutter style={{ margin: "20px 0px" }} sm={6} md={4} lg={3} xl2={1}>
                  <GridItem span={12} rowSpan={5}><DashboardLineChart data={prowJobFailuerMetrics}></DashboardLineChart></GridItem>
                </Grid>
                }

              </div>
              }
            </React.Fragment>
          </PageSection>
        </DrawerContent>
      </Drawer>
    </React.Fragment>

  )
}

export { Reports };