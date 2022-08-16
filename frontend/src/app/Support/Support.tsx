import React, { useEffect, useState, useLayoutEffect, useRef } from 'react';
import { CubesIcon, ExclamationCircleIcon, OkIcon, HelpIcon } from '@patternfly/react-icons';
import {
  PageSection,PageSectionVariants,
  EmptyState,
  EmptyStateVariant,
  EmptyStateIcon,
  EmptyStateBody,
  DataList, DataListItem, DataListItemRow, DataListItemCells, DataListCell,
  Title, TitleSizes,
  Alert, AlertGroup, AlertActionCloseButton
} from '@patternfly/react-core';
import { Toolbar, ToolbarItem, ToolbarContent } from '@patternfly/react-core';
import { Button } from '@patternfly/react-core';
import { Select, SelectOption, SelectVariant } from '@patternfly/react-core';
import { getAllRepositoriesWithOrgs, getLatestProwJob, getProwJobStatisticsMOCK } from '@app/utils/APIService';
import { Card, CardTitle, CardBody } from '@patternfly/react-core';
import { Chart, ChartAxis, ChartBar, ChartLine, ChartGroup, ChartVoronoiContainer } from '@patternfly/react-charts';
import { SimpleList, SimpleListItem } from '@patternfly/react-core';
import { Grid, GridItem } from '@patternfly/react-core';

export interface JobsStatistics {
  repository_name: string;
  type: string;
  git_org: string;
  jobs: (JobsEntity)[] ;
}

export interface JobsEntity {
  name: string;
  metrics: (MetricsEntity)[];
  summary: (MetricsSummary);
}

export interface MetricsSummary {
  success_rate_avg: string;
  failure_rate_avg: string;
  ci_failed_rate_avg: string;
  date_from: string;
  date_to: string;
}

export interface MetricsEntity {
  success_rate: string;
  failure_rate: string;
  ci_failed_rate: string;
  date: string;
}

type SimpleListDemoProps = {
  data: any,
  onSelection: (value) => void
};

const SimpleListDemo = ({data, onSelection}:SimpleListDemoProps) => {
  const onSelect = (selectedItem, selectedItemProps) => {
    onSelection(selectedItemProps["data-index"])
  }

  const items = data.map((job) => <SimpleListItem className="" key={job.index} data-index={job.index} isActive={job.index==0}> {job.job} </SimpleListItem>);

  return (
      <Card style={{width: "100%", height: "100%", fontSize: "1rem"}}>
        <CardTitle>Jobs</CardTitle>
        <CardBody>
          <SimpleList onSelect={onSelect} aria-label="Simple List Example">
            {items}
          </SimpleList>
        </CardBody>
      </Card>
  )
};

const BasicWithRightAlignedLegend = ({data, jobIndex}: {data:JobsStatistics, jobIndex:number}) => {
  const ref = useRef<HTMLDivElement>(null);

  const [width, setWidth] = useState(0);
  const [height, setHeight] = useState(0);

  useLayoutEffect(() => {
    if(ref.current){
      setWidth(ref.current.offsetWidth-20);
      setHeight(ref.current.offsetHeight-50);
    }
  }, []);

  type MyType = {
    x: string;
    y: number;
    name: string;
  }

  let beautifiedData: { [key: string]: MyType[] } = {
    "SUCCESS_RATE_INDEX":[],
    "FAILURE_RATE_INDEX":[],
    "CI_FAILED_RATE_INDEX":[],
    "SUCCESS_RATE_AVG_INDEX":[],
    "FAILURE_RATE_AVG_INDEX":[],
    "CI_FAILED_RATE_AVG_INDEX":[],
  };

  data.jobs[jobIndex].metrics.map(metric => {
    beautifiedData["SUCCESS_RATE_INDEX"].push({name: 'success_rate', x: new Date(metric.date).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: +metric.success_rate})
    beautifiedData["FAILURE_RATE_INDEX"].push({name: 'failure_rate', x: new Date(metric.date).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: +metric.failure_rate})
    beautifiedData["CI_FAILED_RATE_INDEX"].push({name: 'ci_failed_rate', x: new Date(metric.date).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: +metric.ci_failed_rate})
  });

  beautifiedData["SUCCESS_RATE_AVG_INDEX"] = [
    {name: 'success_rate_avg', x: new Date(data.jobs[jobIndex].summary.date_from).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: +data.jobs[jobIndex].summary.success_rate_avg},
    {name: 'success_rate_avg', x: new Date(data.jobs[jobIndex].summary.date_to).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: +data.jobs[jobIndex].summary.success_rate_avg}
  ]

  beautifiedData["FAILURE_RATE_AVG_INDEX"] = [
    {name: 'failure_rate_avg', x: new Date(data.jobs[jobIndex].summary.date_from).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: +data.jobs[jobIndex].summary.failure_rate_avg},
    {name: 'failure_rate_Avg', x: new Date(data.jobs[jobIndex].summary.date_to).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: +data.jobs[jobIndex].summary.failure_rate_avg}
  ]

  beautifiedData["CI_FAILED_RATE_AVG_INDEX"] = [
    {name: 'ci_failed_rate_avg', x: new Date(data.jobs[jobIndex].summary.date_from).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: +data.jobs[jobIndex].summary.ci_failed_rate_avg},
    {name: 'ci_failed_rate_avg', x: new Date(data.jobs[jobIndex].summary.date_to).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: +data.jobs[jobIndex].summary.ci_failed_rate_avg}
  ]


  return (
    <div style={{ height: '100%', width: '100%', minHeight: "600px"}} className={"pf-c-card"} ref={ref}>
      <div style={{ height: height+'px', width: width+'px', background: "white" }}>
        <Chart
          ariaDesc="Average number of pets"
          ariaTitle="Bar chart example"
          containerComponent={<ChartVoronoiContainer labels={({ datum }) => `${datum.name}: ${datum.y}`} constrainToVisibleArea />}
          domain={{y: [-10,110]}}
          legendData={[
            { name: 'success_rate', symbol: {fill: 'rgba(30, 79, 24, 0.5)'} }, 
            { name: 'failure_rate', symbol: {fill: 'rgba(163, 0, 0, 0.5)'}  }, 
            { name: 'ci_failed_rate', symbol: {fill: 'rgba(240, 171, 0, 0.5)'} }, 
            { name: 'success_rate_avg', symbol: {fill: 'rgba(30, 79, 24, 0.2)'} }, 
            { name: 'failure_rate_avg', symbol: {fill: 'rgba(163, 0, 0, 0.2)'} }, 
            { name: 'ci_failed_rate_avg', symbol: {fill: 'rgba(240, 171, 0, 0.2)'} },
          ]}
          legendOrientation="horizontal"
          legendPosition="bottom"
          height={height}
          padding={{
            bottom: 80,
            left: 80,
            right: 80,
            top: 80
          }}
          width={width}
        >
          <ChartAxis></ChartAxis>
          <ChartAxis dependentAxis showGrid />
          <ChartGroup colorScale={['rgba(30, 79, 24, 0.5)', 'rgba(163, 0, 0, 0.5)', 'rgba(240, 171, 0, 0.5)', 'rgba(30, 79, 24, 0.2)', 'rgba(163, 0, 0, 0.2)', 'rgba(240, 171, 0, 0.2)']}>
            <ChartLine data={beautifiedData["SUCCESS_RATE_INDEX"]}/>
            <ChartLine data={beautifiedData["FAILURE_RATE_INDEX"]}/>
            <ChartLine data={beautifiedData["CI_FAILED_RATE_INDEX"]}/>
            <ChartLine data={beautifiedData["SUCCESS_RATE_AVG_INDEX"]} style={{data: { strokeDasharray: '5', strokeWidth: '5'}}}/>
            <ChartLine data={beautifiedData["FAILURE_RATE_AVG_INDEX"]} style={{data: { strokeDasharray: '10', strokeWidth: '5'}}}/>
            <ChartLine data={beautifiedData["CI_FAILED_RATE_AVG_INDEX"]} style={{data: { strokeDasharray: '15', strokeWidth: '5'}}}/>
            </ChartGroup>
        </Chart>
      </div>
    </div>
  )
};

type CardProps = {
  cardType?: 'default' | 'danger' | 'success' | 'warning' | 'primary' | 'help';
  title: string;
  body: string;
};

type InfoCardProp = {
  title: string;
  value: string;
}

const InfoCard = ({data}: {data:InfoCardProp[]}) => {

  return (
    <Card style={{width: "100%", height: "100%", fontSize: "1rem"}}>
      <CardBody>
        { 
          data.map(function(value, index){
            return ( <div style={{marginTop: "5px"}}>
              <div><strong>{value.title}</strong></div>
              <div>{value.value}</div>
              </div> 
            )
          })
        }
      </CardBody>
    </Card>
  )
};

const CardWithNoFooter = ({cardType, title, body}:CardProps) => {
  const cardStyle = new Map();
  cardStyle.set('title-danger', {color: "#A30000", fontWeight: "semibold", fontSize: "0.8em"});
  cardStyle.set('title-success', {color: "#1E4F18", fontWeight: "semibold", fontSize: "0.8em"});
  cardStyle.set('title-warning', {color: "#F0AB00", fontWeight: "semibold", fontSize: "0.8em"});
  cardStyle.set('title-default', {color: "black", fontWeight: "semibold", fontSize: "0.8em"});
  cardStyle.set('title-help', {color: "grey", fontWeight: "semibold", fontSize: "1em"});
  cardStyle.set('title-primary', {color: "#0066CC", fontWeight: "semibold", fontSize: "0.8em"});
  cardStyle.set('body-danger', {color: "#A30000", fontWeight: "bold", fontSize: "2em", textAlign: "center"});
  cardStyle.set('body-success', {color: "#1E4F18", fontWeight: "bold", fontSize: "2em", textAlign: "center"});
  cardStyle.set('body-warning', {color: "#F0AB00", fontWeight: "bold", fontSize: "2em", textAlign: "center"});
  cardStyle.set('body-default', {color: "black", fontWeight: "bold", fontSize: "1.8em", textAlign: "center"});
  cardStyle.set('body-help', {color: "grey", fontWeight: "normal", fontSize: "0.8em", textAlign: "left"});
  cardStyle.set('body-primary', {color: "#0066CC", fontWeight: "bold", fontSize: "2em", textAlign: "center"});

  return (
    <Card style={{width: "100%", height: "100%"}}>
      <CardTitle style={cardStyle.get("title-"+cardType)}>
        {cardType == 'help' && <HelpIcon style={{marginRight: "5px", fontSize: "1.1em", fontWeight: "bold", verticalAlign: "middle"}}></HelpIcon>}  
        {title}
      </CardTitle>
      <CardBody style={cardStyle.get("body-"+cardType)}>
        {body}
        {cardType == 'danger' && <ExclamationCircleIcon style={{fontSize: "1.2rem", margin: "0 5px"}}></ExclamationCircleIcon>}
        {cardType == 'success' && <OkIcon style={{fontSize: "1.2rem", margin: "0 5px"}}></OkIcon>}
        </CardBody>
    </Card>
  )
};

// eslint-disable-next-line prefer-const
let Support = () => {
  const LoremIpsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat"
  const [repositories, setRepositories] = useState<{repoName: string, organization: string, isPlaceholder?: boolean}[]>([]);
  const [repoName, setRepoName] = useState("infra-deployments");
  const [repoOrg, setRepoOrg] = useState("redhat-appstudio");
  const [jobType, setjobType] = useState("periodic");
  const [jobTypeToggle, setjobTypeToggle] = useState(false);
  const [repoNameToggle, setRepoNameToggle] = useState(false);
  const [prowVisible, setProwVisible] = useState(false)
  const [buttonDisabled, setbuttonDisabled] = useState(true);
  const [alerts, setAlerts] = React.useState<React.ReactNode[]>([]);
  const [prowJobs, setprowJobs] = useState([])
  const [selectedJob, setSelectedJob] = useState(0)
  const [prowJobsStats, setprowJobsStats] = useState<JobsStatistics|null>(null);

  const setRepoNameOnChange = (event, selection, isPlaceholder) => { 
    if (isPlaceholder){
      setRepoName("");
      setRepoOrg("");
      setRepoNameToggle(false);
    }
    else {
      setRepoName(repositories[selection].repoName); 
      setRepoOrg(repositories[selection].organization); 
      setRepoNameToggle(false) 
    }
  };

  const setjobTypeOnChange  = (event, selection, isPlaceholder) => { 
    if (isPlaceholder){
      setjobType("");
      setjobTypeToggle(false);
    }
    else{
      setjobType(selection); 
      setjobTypeToggle(false);
    }
  };

  const validateGetProwJob = () => {
    if(repoName != "" && repoOrg != "" && jobType != ""){
      setbuttonDisabled(false)
    }
    else{
      setbuttonDisabled(true)
    }
  }

  useEffect(() => {
    validateGetProwJob();
  }, [repoName, jobType]);

  const getProwJob = async () => {
    setProwVisible(true)
    try {
      let data = await getLatestProwJob(repoName, repoOrg, jobType)
      setprowJobs(data)

      let stats = await getProwJobStatisticsMOCK(repoName, repoOrg, jobType)
      setprowJobsStats(stats)
      
    }
    catch {
      setProwVisible(false);
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
    }
  }

  const clearProwJob = () => {
    setProwVisible(false);
    setRepoName("");
    setRepoOrg("");
    setjobType("");
  }

  const clearRepo = () => {
    setRepoName("")
    setRepoOrg("")
    setRepoNameToggle(false) 
  }

  useEffect( () => {
    getAllRepositoriesWithOrgs()
    .then((data:any) => {
      data.unshift({repoName: "Select a repository", organization: "", isPlaceholder: true})
      setRepositories(data)
    })
  }, []);

  let jobTypes = [
    <SelectOption key={0} value="periodic"/>,
    <SelectOption key={1} value="presubmit"/>,
    <SelectOption key={2} value="postsubmit"/>,
  ]

  let statusColorMap = new Map<string, string>([
    ["skipped", "lightgrey"],
    ["passed", "darkgreen"],
    ["failed", "darkred"]
  ]);

  let jobNames = prowJobsStats?.jobs != null ? prowJobsStats.jobs.map(function(job, index){ return {"job":job.name, "index":index} }) : []

  return (
    
    <React.Fragment>
      {/* page title bar */}
      <PageSection variant={PageSectionVariants.light}>
        <Title headingLevel="h3" size={TitleSizes['2xl']}>
          Tests Reports
        </Title>
      </PageSection>
      {/* main content  */}
      <PageSection>
        {/* alertGroup will show toast notification (on the top left) when an error occurs */}
        <AlertGroup isToast isLiveRegion> {alerts} </AlertGroup>
        {/* the following toolbar will contain the form (dropdowns and button) to request data to the server */}
        <Toolbar style={{ width:  prowVisible ? '100%' : '100%', margin: prowVisible ? 'auto' : '0 auto' }}>
          <ToolbarContent style={{textAlign: 'center'}}>
            <ToolbarItem  style={{minWidth: "20%", maxWidth: "40%"}}>
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
            <ToolbarItem style={{minWidth: "20%", maxWidth: "40%"}}>
              <Select placeholderText="Filter by status/vendor" isOpen={jobTypeToggle} onToggle={setjobTypeToggle} selections={jobType} onSelect={setjobTypeOnChange} aria-label="Select Input">
                {jobTypes}
              </Select>
            </ToolbarItem>
            <ToolbarItem >      
              <Button variant="primary" isDisabled={buttonDisabled} onClick={getProwJob}>Get Latest Test Report</Button>
            </ToolbarItem>
            <ToolbarItem >      
              <Button variant="link" onClick={clearProwJob}>Clear</Button>
            </ToolbarItem>
          </ToolbarContent>
        </Toolbar>
        {/* if the server has not provided any data or if the clear button is clicked or if the page is in its initial state, this empty placeholder will be shown */}
        {!prowVisible && <EmptyState variant={EmptyStateVariant.xl}>
          <EmptyStateIcon icon={CubesIcon}/>
            <Title headingLevel="h1" size="lg">
              No job selected yet.
              </Title>
            <EmptyStateBody>
              Please select a repository and an organization to see the last job's details
            </EmptyStateBody>
          </EmptyState>
        }
        {/* this section will show statistics and details about job and suites */}
        <React.Fragment>
          { prowVisible && <div style={{marginTop: '20px'}}>
            {/* this section will show the job's chart over time and last execution stats */}
            { prowJobsStats !== null && <Grid hasGutter style={{margin: "20px 0px"}} sm={6} md={4} lg={3} xl2={1}>
              <GridItem span={3}><InfoCard data={[{title: "Repository", value: prowJobsStats.repository_name}, {title: "Organization", value: prowJobsStats.git_org}]}></InfoCard></GridItem>
              <GridItem span={2}><CardWithNoFooter cardType={'danger'} title="avg of ci failures" body={prowJobsStats?.jobs != null ? prowJobsStats.jobs[selectedJob].summary.ci_failed_rate_avg +"%" : "-"}></CardWithNoFooter></GridItem>
              <GridItem span={2}><CardWithNoFooter cardType={'danger'} title="avg of failures" body={prowJobsStats?.jobs != null ? prowJobsStats.jobs[selectedJob].summary.failure_rate_avg +"%" : "-"}></CardWithNoFooter></GridItem>
              <GridItem span={2}><CardWithNoFooter cardType={'success'} title="avg of passed tests" body={prowJobsStats?.jobs != null ? prowJobsStats.jobs[selectedJob].summary.success_rate_avg +"%" : "-"}></CardWithNoFooter></GridItem>
              <GridItem span={3}><CardWithNoFooter cardType={'default'} title="Time Range" body={prowJobsStats?.jobs != null ? new Date(prowJobsStats.jobs[selectedJob].summary.date_from.split(" ")[0]).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }) + " - " + new Date(prowJobsStats.jobs[selectedJob].summary.date_to.split(" ")[0]).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }): "-"}></CardWithNoFooter></GridItem>

              <GridItem span={4} rowSpan={4}><SimpleListDemo data={jobNames} onSelection={(value)=>{setSelectedJob(value)}}></SimpleListDemo></GridItem>
              <GridItem span={8} rowSpan={5}><BasicWithRightAlignedLegend data={prowJobsStats} jobIndex={selectedJob}></BasicWithRightAlignedLegend></GridItem>
              <GridItem span={4} rowSpan={1}><CardWithNoFooter cardType={'help'} title="About this dashboard" body={LoremIpsum}></CardWithNoFooter></GridItem>

            </Grid> 
            }
            {/* TDB: if job has also suite details, here is the place to show it */}
            { false && <DataList aria-label="Simple data list example" style={{marginTop: '20px'}}>
              <DataListItem aria-labelledby="simple-item1">
                <DataListItemRow key="000" style={{fontWeight: 'bold', borderBottom: "1px solid lightgrey"}}>
                  <DataListItemCells
                    dataListCells={[
                      <DataListCell width={4} key={2}>Name</DataListCell>,
                      <DataListCell width={1} key={3}>Status</DataListCell>,
                      <DataListCell width={1} key={4}>Time elapsed</DataListCell>,
                    ]}
                  />
                </DataListItemRow>
                {prowJobs.map(function(value, index){
                  return <DataListItemRow key={index}>
                        <DataListItemCells
                          dataListCells={[
                            <DataListCell width={4} key={index+"-2"}>{value['name']}</DataListCell>,
                            <DataListCell width={1} key={index+"-3"} style={{fontWeight: "bold", textTransform: "uppercase", color : statusColorMap.get(value["status"])}}>
                              {value['status']}
                            </DataListCell>,
                            <DataListCell width={1} key={index+"-4"}>{value['time']}</DataListCell>,
                          ]}
                        />
                    </DataListItemRow>
                })}
              </DataListItem>
            </DataList>
            }
          </div> 
          }
        </React.Fragment>
      </PageSection>
    </React.Fragment>

)}

export { Support };
