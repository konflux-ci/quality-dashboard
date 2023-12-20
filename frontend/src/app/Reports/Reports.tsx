import React, { useEffect, useState } from 'react';
import { CopyIcon, CubesIcon, ExclamationCircleIcon, OkIcon, InfoCircleIcon } from '@patternfly/react-icons';
import {
  PageSection, PageSectionVariants,
  EmptyState,
  EmptyStateVariant,
  EmptyStateIcon,
  EmptyStateBody,
  Title, TitleSizes,
  Alert, AlertGroup, AlertActionCloseButton,
  Spinner, Breadcrumb, BreadcrumbItem, Tooltip
} from '@patternfly/react-core';
import { Toolbar, ToolbarItem, ToolbarContent } from '@patternfly/react-core';
import { Button } from '@patternfly/react-core';
import { Select, SelectOption, SelectVariant } from '@patternfly/react-core';
import { getAllRepositoriesWithOrgs, getFlakyData, getJobTypes, getProwJobStatistics, getProwJobs, getTeams } from '@app/utils/APIService';
import { Grid, GridItem } from '@patternfly/react-core';
import {
  JobsStatistics,
  DashboardLineChart,
  SimpleListData,
  DashboardLineChartData
} from '@app/utils/sharedComponents';
import { Card, CardTitle, CardBody, CardFooter } from '@patternfly/react-core';
import { Flex, FlexItem } from '@patternfly/react-core';

import { Divider, TextVariants, Text } from '@patternfly/react-core';
import { ReactReduxContext, useSelector } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { formatDate, getRangeDates } from './utils';
import { DateTimeRangePicker } from '../utils/DateTimeRangePicker';
import { Job } from './FailedE2ETests';
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

  const [repositories, setRepositories] = useState<{ repoName: string, organization: string, isPlaceholder?: boolean }[]>([]);
  const [repoName, setRepoName] = useState("");
  const [repoOrg, setRepoOrg] = useState("");
  const [jobType, setjobType] = useState("");
  const [jobTypes, setJobTypes] = useState<string[]>([]);
  const [jobTypeToggle, setjobTypeToggle] = useState(false);
  const [jobNameToggle, setjobNameToggle] = useState(false);
  const [repoNameToggle, setRepoNameToggle] = useState(false);
  const currentTeam = useSelector((state: any) => state.teams.Team);
  const history = useHistory();
  const params = new URLSearchParams(window.location.search);
  const [rangeDateTime, setRangeDateTime] = useState(getRangeDates(10));
  const [isInvalid, setIsInvalid] = useState(false);

  // Called onChange of the repository dropdown element. This set repository name and organization state variables, or clears them when placeholder is selected
  const setRepoNameOnChange = (event, selection, isPlaceholder) => {
    if (isPlaceholder) {
      clearRepo()
    }
    else {
      setRepoName(repositories[selection].repoName);
      setRepoOrg(repositories[selection].organization);
      setRepoNameToggle(false)
      params.set("repository", repositories[selection].repoName)
      params.set("organization", repositories[selection].organization)

      getJobTypes(repositories[selection].repoName, repositories[selection].organization)
        .then((data: any) => {
          setJobTypes(data)
          setjobType("presubmit") // all repos in OpenShift CI have presubmit type job
          params.set("job_type", "presubmit")
          params.set("start", formatDate(rangeDateTime[0]))
          params.set("end", formatDate(rangeDateTime[1]))
          history.push(window.location.pathname + '?' + params.toString());
        });
    }
  };

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
    setSelectedJobObject({});
    setSelectedJob(0);
    setjobNameToggle(false);
  }

  // Reset rangeDateTime
  const clearRangeDateTime = () => {
    setRangeDateTime(getRangeDates(10))
  }

  // Called onChange of the jobType dropdown element. This set repository name and organization state variables, or clears them when placeholder is selected
  const setjobTypeOnChange = (event, selection, isPlaceholder) => {
    if (isPlaceholder) {
      clearJobType()
    }
    else {
      setjobType(selection);
      setjobTypeToggle(false);
      setjobNameToggle(false);
      setIsInvalid(false)
      params.set("job_type", selection)
      history.push(window.location.pathname + '?' + params.toString());
    }
  };

  // Validates that the required variables are not empty; if not, the "get" button is enabled
  const validateGetProwJob = () => {
    if (repositories.find(r => r.organization == repoOrg && r.repoName == repoName)) {
      getJobTypes(repoName, repoOrg)
        .then((data: any) => {
          if (data.find(j => j == jobType)) {
            getProwJob()
          }
        });
    }
  }

  // Triggers automatic validation when state variables change
  useEffect(() => {
    validateGetProwJob();
  }, [repoOrg, repoName, jobType, rangeDateTime]);

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
      const start = params.get("start")
      const end = params.get("end")

      getAllRepositoriesWithOrgs(state.teams.Team, true, rangeDateTime)
        .then((data: any) => {
          let dropDescr = ""
          if (data.length < 1 && (team == state.teams.Team || team == null)) {
            dropDescr = "No Repositories"
            setLoadingState(false)
            history.push('/reports/test?team=' + currentTeam)
          }
          else { dropDescr = "Select a repository" }

          if (data.length > 0 && (team == state.teams.Team || team == null)) {
            data.unshift({ repoName: dropDescr, organization: "", isPlaceholder: true }) // Adds placeholder at the beginning of the array, so it will be shown first
            setRepositories(data)

            if (repository == null || organization == null || job_type == null || start == null || end == null) { // first click on OpenShift CI or team
              setRepoName(data[1].repoName)
              setRepoOrg(data[1].organization)
              setjobType("presubmit") // all repos in OpenShift CI have presubmit type job

              getJobTypes(data[1].repoName, data[1].organization)
                .then((data: any) => {
                  setJobTypes(data)
                });

              const start_date = formatDate(rangeDateTime[0])
              const end_date = formatDate(rangeDateTime[1])

              history.push('/reports/test?team=' + currentTeam + '&organization=' + data[1].organization + '&repository=' + data[1].repoName
                + '&job_type=presubmit' + '&start=' + start_date + '&end=' + end_date)

            } else {
              if (validateRepositoryParams(data, repository, organization)) {
                setRepoName(repository)
                setRepoOrg(organization)
                setRangeDateTime([new Date(start), new Date(end)])

                getJobTypes(repository, organization)
                  .then((data: any) => {
                    setJobTypes(data)
                    if (validateParam(data, job_type)) {
                      setjobType(job_type)
                      history.push('/reports/test?team=' + currentTeam + '&organization=' + organization + '&repository=' + repository +
                        '&job_type=' + job_type + '&start=' + start + '&end=' + end)
                    } else {
                      setLoadingState(false)
                      setIsInvalid(true)
                    }
                  });
              } else {
                setLoadingState(false)
                setIsInvalid(true)
              }
            }
          }
        })
    }
  }, [setRepositories, currentTeam]);

  /* 
  ProwJobs logic to populate dashboard
  */

  const [selectedJob, setSelectedJob] = useState(0)
  const [selectedJobObject, setSelectedJobObject] = useState({})

  const [prowJobsStats, setprowJobsStats] = useState<JobsStatistics | null>(null);
  const [prowJobs, setProwJobs] = useState<Job[] | null>(null);

  // Get the prow jobs from API
  const getProwJob = async () => {
    setSelectedJob(0)
    // Hide components and show loading spinner 
    setProwVisible(false)
    setLoadingState(true)
    setNoData(false)
    try {
      // Get statistics and metrics
      const stats = await getProwJobStatistics(repoName, repoOrg, jobType, rangeDateTime)
      // Get jobs
      const prowJobs = await getProwJobs(repoName, repoOrg, rangeDateTime)
      setProwJobs(prowJobs)
      // Set UI for showing data and disable spinner
      setprowJobsStats(stats)
      setLoadingState(false)
      setProwVisible(true)

      if (stats != null) {
        const data = await getFlakyData(currentTeam, stats.jobs[0].name, repoName, rangeDateTime[0].toISOString(), rangeDateTime[1].toISOString(), repoOrg)
        setImpact(data.data.global_impact)
      } else {
        setImpact("")
      }
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

  const ci_html: string = prowJobsStats?.jobs != null ? "https://prow.ci.openshift.org/?repo=" + prowJobsStats?.git_organization + "%2F" + prowJobsStats?.repository_name + "&type=" + prowJobsStats?.type : ''

  // Extract a simple list of jobs from data: this will be used to let users select the job they want to see details for
  const jobNames: SimpleListData[] = prowJobsStats?.jobs != null ? prowJobsStats.jobs.map(function (job, index) { return { "value": job.name + " (Total: " + job.summary.total_jobs + ")", "index": index, "href": ci_html + "&job=" + job.name } }) : []

  // Prepare data for the line chart
  const beautifiedData: DashboardLineChartData = {
    "SUCCESS_RATE_INDEX": { data: [] },
    "FAILURE_RATE_INDEX": { data: [] },
    "CI_FAILED_RATE_INDEX": { data: [] },
    "SUCCESS_RATE_AVG_INDEX": { data: [] },
    "FAILURE_RATE_AVG_INDEX": { data: [] },
    "CI_FAILED_RATE_AVG_INDEX": { data: [] },
  };

  if (prowJobsStats) {
    prowJobsStats.jobs[selectedJob].metrics.map(metric => {
      beautifiedData["SUCCESS_RATE_INDEX"].data.push({ name: 'success_rate', x: new Date(metric.date).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: + (metric.total_jobs > 0 ? (100 * metric.success_count / metric.total_jobs).toFixed(2) : 0) })
      beautifiedData["FAILURE_RATE_INDEX"].data.push({ name: 'failure_rate', x: new Date(metric.date).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: +  (metric.total_jobs > 0 ? (100 * metric.failure_count / metric.total_jobs).toFixed(2) : 0) })
      beautifiedData["CI_FAILED_RATE_INDEX"].data.push({ name: 'ci_failed_rate', x: new Date(metric.date).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: + (metric.total_jobs > 0 ? (100 * metric.not_scheduled_count / metric.total_jobs).toFixed(2) : 0) })
    });

    beautifiedData["SUCCESS_RATE_INDEX"].style = { data: { stroke: "rgba(30, 79, 24, 0.9)", strokeWidth: 2 } }
    beautifiedData["FAILURE_RATE_INDEX"].style = { data: { stroke: "rgba(163, 0, 0, 0.9)", strokeWidth: 2 } }
    beautifiedData["CI_FAILED_RATE_INDEX"].style = { data: { stroke: "rgba(240, 171, 0, 0.9)", strokeWidth: 2 } }

    beautifiedData["SUCCESS_RATE_AVG_INDEX"].data = [
      { name: 'success_rate_avg', x: new Date(prowJobsStats.jobs[selectedJob].summary.date_from).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: + (100 * prowJobsStats.jobs[selectedJob].summary.success_count / prowJobsStats.jobs[selectedJob].summary.total_jobs).toFixed(2) },
      { name: 'success_rate_avg', x: new Date(prowJobsStats.jobs[selectedJob].summary.date_to).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: + (100 * prowJobsStats.jobs[selectedJob].summary.success_count / prowJobsStats.jobs[selectedJob].summary.total_jobs).toFixed(2) }
    ]
    beautifiedData["SUCCESS_RATE_AVG_INDEX"].style = { data: { stroke: "rgba(30, 79, 24, 0.3)", strokeDasharray: 10, strokeWidth: 5 } }

    beautifiedData["FAILURE_RATE_AVG_INDEX"].data = [
      { name: 'failure_rate_avg', x: new Date(prowJobsStats.jobs[selectedJob].summary.date_from).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: + (100 * prowJobsStats.jobs[selectedJob].summary.failure_count / prowJobsStats.jobs[selectedJob].summary.total_jobs).toFixed(2) },
      { name: 'failure_rate_avg', x: new Date(prowJobsStats.jobs[selectedJob].summary.date_to).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: + (100 * prowJobsStats.jobs[selectedJob].summary.failure_count / prowJobsStats.jobs[selectedJob].summary.total_jobs).toFixed(2) }
    ]
    beautifiedData["FAILURE_RATE_AVG_INDEX"].style = { data: { stroke: "rgba(163, 0, 0, 0.3)", strokeDasharray: 10, strokeWidth: 5 } }

    beautifiedData["CI_FAILED_RATE_AVG_INDEX"].data = [
      { name: 'ci_failed_rate_avg', x: new Date(prowJobsStats.jobs[selectedJob].summary.date_from).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: + (100 * prowJobsStats.jobs[selectedJob].summary.not_scheduled_count / prowJobsStats.jobs[selectedJob].summary.total_jobs).toFixed(2) },
      { name: 'ci_failed_rate_avg', x: new Date(prowJobsStats.jobs[selectedJob].summary.date_to).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: + (100 * prowJobsStats.jobs[selectedJob].summary.not_scheduled_count / prowJobsStats.jobs[selectedJob].summary.total_jobs).toFixed(2) }
    ]
    beautifiedData["CI_FAILED_RATE_AVG_INDEX"].style = { data: { stroke: "rgba(240, 171, 0, 0.3)", strokeDasharray: 10, strokeWidth: 5 } }
  }

  function handleChange(event, from, to) {
    setRangeDateTime([from, to])
    params.set("start", formatDate(from))
    params.set("end", formatDate(to))
    history.push(window.location.pathname + '?' + params.toString());
  }

  const start = rangeDateTime[0]
  const end = rangeDateTime[1]

  useEffect(() => {
    if (prowJobsStats != null) {
      getFlakyData(currentTeam, prowJobsStats.jobs[selectedJob].name, repoName, start.toISOString(), end.toISOString(), repoOrg).then(res => { setImpact(res.data.global_impact) })
    } else {
      setImpact("")
    }
  }, [selectedJob]);

  return (

    <React.Fragment>
      {/* page title bar */}
      <Header info="Observe the CI metrics of all your Ci jobs."></Header>
      <div style={{ padding: 15, background: "white", border: "1px solid lightgrey" }}>
        <Breadcrumb>
          <BreadcrumbItem>OpenShift CI</BreadcrumbItem>
          <BreadcrumbItem>Tests Reports</BreadcrumbItem>
        </Breadcrumb>
      </div>
      <PageSection variant={PageSectionVariants.light}>
        <Title headingLevel="h3" size={TitleSizes['2xl']}>
          Tests Reports
          <Button onClick={() => navigator.clipboard.writeText(window.location.href)} variant="link" icon={<CopyIcon />} iconPosition="right">
            Copy link
          </Button>
        </Title>
      </PageSection>

      <PageSection variant={PageSectionVariants.light} style={{ padding: "10px 5px", background: "white" }}>  
        {/* the following toolbar will contain the form (dropdowns and button) to request data to the server */}
        <Toolbar style={{ width: prowVisible ? '100%' : '100%', margin: prowVisible ? 'auto' : '0', padding: 0 }}>
          <ToolbarContent style={{ textAlign: 'left' }}>
            <ToolbarItem style={{ minWidth: "20%", maxWidth: "40%" }}>
              <span id="typeahead-select" hidden>
                Select a state
              </span>
              <Select
                variant={SelectVariant.typeahead}
                typeAheadAriaLabel="Select a repository"
                isOpen={repoNameToggle}
                onToggle={setRepoNameToggle}
                selections={repoName}
                onSelect={setRepoNameOnChange}
                onClear={clearRepo}
                aria-labelledby="typeahead-select"
                placeholderText="Select a repository"
              >
                {repositories.map((value, index) => (
                  <SelectOption key={index} value={index} description={value.organization} isDisabled={value.isPlaceholder}>{value.repoName}</SelectOption>
                ))}
              </Select>
            </ToolbarItem>
            <ToolbarItem style={{ minWidth: "20%", maxWidth: "40%" }}>
              <Select placeholderText="Filter by status/vendor" isOpen={jobTypeToggle} onToggle={setjobTypeToggle} selections={jobType} onSelect={setjobTypeOnChange} aria-label="Select Input">
                {jobTypes.map((value, index) => (
                  <SelectOption key={index} value={value}>{value}</SelectOption>
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
            {repoName && jobType &&
              <ToolbarItem style={{ minWidth: "20%", maxWidth: "40%" }}>
              
                <Select isOpen={jobNameToggle} onToggle={setjobNameToggle} selections={jobNames[selectedJob] ? jobNames[selectedJob].value : "Select a job"} onSelect={(e, value) => { setSelectedJob(value.index); setSelectedJobObject(value); setjobNameToggle(false) }} aria-label="Select Input">
                  { jobNames.map((value, index) => (
                    <SelectOption key={index} isPlaceholder={false} value={value}>{value.value}</SelectOption>
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
              <Grid hasGutter sm={6} md={4} lg={3} xl2={1}>
                <GridItem span={12} rowSpan={1}>
                  <Card>
                    <CardTitle component="h1"><Title headingLevel="h1"><InfoCircleIcon style={{ fontSize: "2rem"}}></InfoCircleIcon></Title></CardTitle>
                    <CardBody>
                      <Text component={TextVariants.p}>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Qrudqwb, vjdskbnvckjvdsnvckjdsn.</Text>
                    </CardBody>
                  </Card>
                </GridItem>
                <GridItem span={5} rowSpan={1}>
                  <Card>
                    <CardTitle component="h1"><Title headingLevel="h1">Jobs Executed</Title></CardTitle>
                    <CardBody>
                      <Text component={TextVariants.p}>Number of jobs executed in the selected time range, with success and failure rate. </Text>
                      <br></br>
                      <Flex className="example-border" justifyContent={{ default: 'justifyContentSpaceBetween' }}>
                        <FlexItem></FlexItem>
                        <FlexItem>
                          <Card style={{border: 'none', boxShadow: "none", textAlign: "center"}}>
                            <CardTitle component="h1">
                              <Title headingLevel="h1"  size={TitleSizes['4xl']}>
                                {prowJobsStats?.jobs != null ? prowJobsStats.jobs[selectedJob].summary.total_jobs : "-"}
                              </Title>
                            </CardTitle>
                            <CardBody>
                              <Title headingLevel="h1">Total</Title>
                            </CardBody>
                            <CardFooter style={{color: "black"}}>
                              Job Runs 
                              <Tooltip content={ <div>Total number of jobs executed in the selected time range</div> }>
                                <InfoCircleIcon style={{ fontSize: "1.2rem", margin: "0 5px" }}></InfoCircleIcon>
                              </Tooltip>
                            </CardFooter>
                          </Card>
                        </FlexItem>
                        <Divider orientation={{ default: 'vertical' }} />
                        <FlexItem>
                          <Card style={{border: 'none', boxShadow: "none", textAlign: "center", color: "#1E4F18"}}>
                            <CardTitle component="h1">
                              <Title headingLevel="h1" size={TitleSizes['4xl']}>
                                {prowJobsStats?.jobs != null ? (100 * prowJobsStats.jobs[selectedJob].summary.success_count / prowJobsStats.jobs[selectedJob].summary.total_jobs).toFixed(2) + "%" : "-"}
                              </Title>
                            </CardTitle>
                            <CardBody>
                              <Title headingLevel="h1">
                                {prowJobsStats?.jobs != null ? ( prowJobsStats.jobs[selectedJob].summary.success_count +"/" +prowJobsStats.jobs[selectedJob].summary.total_jobs ) : "-"}
                                <OkIcon style={{ fontSize: "1.2rem", margin: "0 5px" }}></OkIcon>
                              </Title>
                            </CardBody>
                            <CardFooter style={{color: "black"}}>
                              Completed Jobs 
                              <Tooltip content={ <div>Count and percentage of jobs that completed successfully</div> }>
                                <InfoCircleIcon style={{ fontSize: "1.2rem", margin: "0 5px" }}></InfoCircleIcon>
                              </Tooltip>
                            </CardFooter>
                          </Card>
                        </FlexItem>
                        <Divider orientation={{ default: 'vertical' }} />
                        <FlexItem>
                          <Card style={{border: 'none', boxShadow: "none", textAlign: "center", color: "#A30000"}}>
                            <CardTitle component="h1">
                              <Title headingLevel="h1" size={TitleSizes['4xl']}>
                                {prowJobsStats?.jobs != null ? (100 * prowJobsStats.jobs[selectedJob].summary.failure_count / prowJobsStats.jobs[selectedJob].summary.total_jobs).toFixed(2) + "%" : "-"}
                              </Title>
                            </CardTitle>
                            <CardBody>
                              <Title headingLevel="h1">
                                {prowJobsStats?.jobs != null ? ( prowJobsStats.jobs[selectedJob].summary.failure_count +"/" +prowJobsStats.jobs[selectedJob].summary.total_jobs ) : "-"}
                                <ExclamationCircleIcon style={{ fontSize: "1.2rem", margin: "0 5px" }}></ExclamationCircleIcon>
                              </Title>
                            </CardBody>
                            <CardFooter style={{color: "black"}}>
                              Failed Jobs 
                              <Tooltip content={ <div>Count and percentage of jobs that failed. See the Failures card to understand the most relevant reasons of failure.</div> }>
                                <InfoCircleIcon style={{ fontSize: "1.2rem", margin: "0 5px" }}></InfoCircleIcon>
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
                  <Card>
                    <CardTitle component="h1"><Title headingLevel="h1">Failures</Title></CardTitle>
                    <CardBody>
                      <Text component={TextVariants.p}>The percentage of failures grouped by most common reasons, considering the failed jobs in the selected time range.</Text>
                      <br></br>
                      <Flex className="example-border" justifyContent={{ default: 'justifyContentSpaceBetween' }}>
                        <FlexItem></FlexItem>
                        <FlexItem>
                          <Card style={{border: 'none', boxShadow: "none", textAlign: "center", color: "#A30000"}}>
                            <CardTitle component="h1">
                              <Title headingLevel="h1" size={TitleSizes['4xl']}>
                                {prowJobsStats?.jobs != null ? (100 * prowJobsStats.jobs[selectedJob].summary.not_scheduled_count / prowJobsStats.jobs[selectedJob].summary.failure_count).toFixed(2) + "%" : "-"}
                              </Title>
                            </CardTitle>
                            <CardBody>
                              <Title headingLevel="h1">
                                {prowJobsStats?.jobs != null ? ( prowJobsStats.jobs[selectedJob].summary.not_scheduled_count +"/" +prowJobsStats.jobs[selectedJob].summary.failure_count ) : "-"}
                              </Title>
                            </CardBody>
                            <CardFooter style={{color: "black"}}>
                              CI Fail 
                              <Tooltip content={ <div>Count and percentage of jobs that failed due to CI infrastructure, considering the failing jobs in the selected period of time.</div> }>
                                <InfoCircleIcon style={{ fontSize: "1.2rem", margin: "0 5px" }}></InfoCircleIcon>
                              </Tooltip>
                            </CardFooter>
                          </Card>
                        </FlexItem>
                        <Divider orientation={{ default: 'vertical' }} />
                        <FlexItem>
                          <Card style={{border: 'none', boxShadow: "none", textAlign: "center", color: "#A30000"}}>
                            <CardTitle component="h1">
                              <Title headingLevel="h1"  size={TitleSizes['4xl']}>
                                {impact != "" ? impact + "%" : "-"}
                              </Title>
                            </CardTitle>
                            <CardBody>
                              <Title headingLevel="h1">7/28</Title>
                            </CardBody>
                            <CardFooter style={{color: "black"}}> 
                              {prowJobsStats.jobs[selectedJob].summary.failure_by_e2e_tests_count != 0 ? <a href={'/home/flaky?team=' + currentTeam + '&repository=' + repoName + '&job=' + prowJobsStats.jobs[selectedJob].name + '&start=' + start.toISOString() + '&end=' + end.toISOString()}>Flaky Tests (see more)</a> : <div></div>}
                              <Tooltip content={ <div>Count and percentage of jobs that failed due to flaky tests, considering the failing jobs in the selected period of time. See the flaky tests page to see more details.</div> }>
                                <InfoCircleIcon style={{ fontSize: "1.2rem", margin: "0 5px" }}></InfoCircleIcon>
                              </Tooltip>
                            </CardFooter>
                          </Card>
                        </FlexItem>
                        <Divider orientation={{ default: 'vertical' }} />
                        <FlexItem>
                          <Card style={{border: 'none', boxShadow: "none", textAlign: "center", color: "#A30000"}}>
                          <CardTitle component="h1">
                              <Title headingLevel="h1" size={TitleSizes['4xl']}>
                                {prowJobsStats?.jobs != null ? (100 * prowJobsStats.jobs[selectedJob].summary.failure_by_build_errors_count / prowJobsStats.jobs[selectedJob].summary.failure_count).toFixed(2) + "%" : "-"}
                              </Title>
                            </CardTitle>
                            <CardBody>
                              <Title headingLevel="h1">
                                {prowJobsStats?.jobs != null ? ( prowJobsStats.jobs[selectedJob].summary.failure_by_build_errors_count +"/" +prowJobsStats.jobs[selectedJob].summary.failure_count ) : "-"}
                              </Title>
                            </CardBody>
                            <CardFooter style={{color: "black"}}>
                              Other reasons
                              <Tooltip content={ <div>Count and percentage of jobs that failed due to other reasons (like RHTAP installation failures, cluster build errors, etc.) considering the failing jobs in the selected period of time.</div> }>
                                <InfoCircleIcon style={{ fontSize: "1.2rem", margin: "0 5px" }}></InfoCircleIcon>
                              </Tooltip>
                            </CardFooter>
                          </Card>
                        </FlexItem>
                        <Divider orientation={{ default: 'vertical' }} />
                        <FlexItem>
                          <Card style={{border: 'none', boxShadow: "none", textAlign: "center", color: "#A30000"}}>
                            <CardTitle component="h1"><Title headingLevel="h1" size={TitleSizes['4xl']}>6.18%</Title></CardTitle>
                            <CardBody>
                              <Title headingLevel="h1"> 6/28 </Title>
                            </CardBody>
                            <CardFooter style={{color: "black"}}>
                              CI Fail 
                              <Tooltip content={ <div>Count and percentage of jobs that failed due to external services outage (like Github, Quay.io, etc.), considering the failing jobs in the selected period of time.</div> }>
                                <InfoCircleIcon style={{ fontSize: "1.2rem", margin: "0 5px" }}></InfoCircleIcon>
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
            {prowJobsStats !== null && <Grid hasGutter style={{ margin: "20px 0px" }} sm={6} md={4} lg={3} xl2={1}>      
              <GridItem span={12} rowSpan={5}><DashboardLineChart data={beautifiedData}></DashboardLineChart></GridItem>
            </Grid>
            }

          </div>
          }
        </React.Fragment>
      </PageSection>
    </React.Fragment>

  )
}

export { Reports };