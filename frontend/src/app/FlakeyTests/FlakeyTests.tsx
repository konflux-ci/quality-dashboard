import React, { useState, useRef, useLayoutEffect } from 'react';
import {
  TableComposable,
  Thead,
  Tr,
  Th,
  Tbody,
  Td,
  InnerScrollContainer,
  ExpandableRowContent,
  ThProps
} from '@patternfly/react-table';
import { Header } from '@app/utils/Header';
import { TextContent, Text, TextVariants } from '@patternfly/react-core';
import { PageSection, PageSectionVariants, Title, TitleSizes } from '@patternfly/react-core';
import { Dropdown, DropdownToggle, DropdownItem } from '@patternfly/react-core';
import { ChartPie } from '@patternfly/react-charts';
import { Chart, ChartAxis, ChartLine, ChartGroup, ChartVoronoiContainer } from '@patternfly/react-charts';
import { Grid, GridItem, } from '@patternfly/react-core';
import { Card, CardTitle, CardBody } from '@patternfly/react-core';
import { Toolbar, ToolbarContent } from '@patternfly/react-core';
import { getRepositoriesWithJobs, getFlakyData, getGlobalImpactData } from '../utils/APIService'
import { useSelector } from 'react-redux';
import { DatePicker, Button } from '@patternfly/react-core';
import { Skeleton } from '@patternfly/react-core';

export const DatePickerMinMax: React.FunctionComponent<{selectedDate: Date | undefined, onChange:(value:string, date:Date, name:string)=>void, name:string}> = ({ selectedDate, onChange, name}) => {
  const onDateChange = (e, value, date) => {
    onChange(value, date, name)
  }
  return <DatePicker name={name} onChange={onDateChange} value={selectedDate?.toDateString()} />;
};

const ImpactChart:React.FunctionComponent<{data, x, y}> = ({data, x, y}) => {
  const ref = useRef<HTMLDivElement>(null);
  const [width, setWidth] = useState(100);
  const [height, setHeight] = useState(100);

  useLayoutEffect(() => {
    if (ref.current && ref.current.offsetWidth>0) {
      setWidth(ref.current.offsetWidth * 0.8 - 10);
      setHeight(ref.current.offsetWidth * 0.8 * 0.4 -20);
    }
  }, []);

  return (
    <div style={{  width: '100%', height: '100%', boxShadow: "none" }} className={"pf-c-card"} ref={ref}>
      <div style={{ height: height + 'px', width: width + 'px', background: "white", boxShadow: "none" }}>
        {
          data && data.length > 0 && <Chart
          ariaDesc="Average number of pets"
          ariaTitle="Bar chart example"
          containerComponent={<ChartVoronoiContainer labels={({ datum }) => `${datum.name}: ${datum.y}`} constrainToVisibleArea />}
          domain={{y: [-1,100]}}
          domainPadding={{ x: [30, 25] }}
          legendOrientation="vertical"
          legendPosition="right"
          height={height}
          width={width}
          name="chart1"
          padding={{
            bottom: 30,
            left: 60,
            right: 0,
            top: 50
          }}
        >
          <ChartAxis showGrid style={{ tickLabels: {angle :0, fontSize: 9}}} />
          <ChartAxis dependentAxis showGrid />
          <ChartGroup>
            <ChartLine data={ data.map( (datum) => { return {"name": datum[x], "x": datum[x], "y": datum[y] ? parseFloat(datum[y]) : 0}  }) }/>
          </ChartGroup>
        </Chart>
        }
        {(!data || data.length == 0) && 
          <div style={{ height: '100%', display: 'flex', alignItems: 'flex-end', justifyContent: 'space-between' }}>
            <Skeleton height="25%" width="15%" screenreaderText="Loading contents" />
            <Skeleton height="33%" width="15%" />
            <Skeleton height="50%" width="15%" />
            <Skeleton height="66%" width="15%" />
            <Skeleton height="75%" width="15%" />
            <Skeleton height="100%" width="15%" />
          </div>
        }
      </div>
    </div>
  )
}

const PieChart:React.FunctionComponent<{data, x, y}> = ({data, x, y}) => {
  const ref = useRef<HTMLDivElement>(null);
  const [width, setWidth] = useState(0);
  const [height, setHeight] = useState(0);

  useLayoutEffect(() => {
    if (ref.current) {
      setWidth(ref.current.offsetWidth * 0.8);
      setHeight(ref.current.offsetWidth * 0.8 * 0.6 );
    }
  }, []);

  return (
   <div style={{  width: '100%', height: '100%', boxShadow: "none" }} className={"pf-c-card"} ref={ref}>
      <div style={{ height: height + 'px', width: width + 'px', background: "white", boxShadow: "none" }}>
        {data && data.length > 0 &&
        <ChartPie
          ariaDesc="Average number of pets"
          ariaTitle="Pie chart example"
          constrainToVisibleArea
          colorScale={["tomato", "orange", "gold", "bisque", "coral", "darkorange", "darksalmon", "salmon", "peachpuff", "papayawhip", "palevioletred", "pink", "red" ]}
          data={ data.map((datum => { return {x: datum[x], y: datum[y]} })) }
          height={height}
          width={width}
          legendData={data.map(datum => { return {name: datum[x]} })}
          labels={({ datum }) => `${datum.x}: ${datum.y}`}
          style={{ labels: {fontSize: '9px'} }}
          legendOrientation="vertical"
          legendPosition="right"
          name="chart1"
          padding={{
            bottom: 30,
            left: 20,
            right: 150,
            top: 50
          }}
        />
      }
      {(!data || data.length == 0) && 
        <div style={{ height: '100%', padding:"5%", display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <Skeleton shape="circle" width="50%" screenreaderText="Loading medium circle contents" />
        </div>
      }
      </div>
    </div>
    )
}

export const DropdownBasic: React.FunctionComponent<{toggles, onSelect, selected, placeholder}> = ({toggles, onSelect, selected, placeholder}) => {
  const [isOpen, setIsOpen] = React.useState(false);

  const onToggle = (isOpen: boolean) => {
    setIsOpen(isOpen);
  };

  const onFocus = () => {
    const element = document.getElementById('toggle-basic');
    element?.focus();
  };

  const onItemSelect = (e) => {
    setIsOpen(false);
    onFocus();
    onSelect(e.target.name)
  };
  
  return (
    <Dropdown
      onSelect={onItemSelect}
      toggle={<DropdownToggle onToggle={onToggle}>{selected == '' ? placeholder : selected}</DropdownToggle>}
      isOpen={isOpen}
      dropdownItems={toggles}
      />
  );
};

interface Flakey {
  status: string;
  test_cases: {
    name: string;
    test_case_impact: number;
    count: number;
    messages: {
      job_id: string;
      job_url: string;
      error_message: string;
      failure_date: string;
    }[]
  }[];
  suite_name: string;
  average_impact: number;
}

interface FlakeyObject {
  global_impact: number;
  git_organization: string;
  repository_name: string;
  job_name: string;
  suites: Flakey[];
}

export const ComposableTableNestedExpandable: React.FunctionComponent<{teams:Flakey[]}> = ({teams}) => {

  const columnNames = {
    name: 'Test Case',
    status: 'Status',
    error_message: 'Failing Jobs',
    count: 'Count',
    suite_name: 'Suite Name',
    job_id: "Job ID",
    job_url: "Job Url",
    failure_date: "Falure Date",
    average_impact: "Impact"
  };

  // Exapndable suites
  const [expandedSuitesNames, setExpandedSuitesNames] = React.useState<string[]>([]);

  const setSuiteExpanded = (suite_name: string, isExpanding = true) => {
    setExpandedSuitesNames(prevExpanded => {
      const otherExpandedSuiteNames = prevExpanded.filter(t => t !== suite_name);
      return isExpanding ? [...otherExpandedSuiteNames, suite_name] : otherExpandedSuiteNames;
    });
  }
  const isSuiteExpanded = (suite_name: string) => expandedSuitesNames.includes(suite_name);

  // Expandable test cases
  const [expandedTestCaseNames, setExpandedTestCaseNames] = React.useState<string[]>([]);

  const setTestCaseExpanded = (suite_name: string, isExpanding = true) => {
    setExpandedTestCaseNames(prevExpanded => {
      const otherExpandedTestCaseNames = prevExpanded.filter(t => t !== suite_name);
      return isExpanding ? [...otherExpandedTestCaseNames, suite_name] : otherExpandedTestCaseNames;
    });
  }
  const isTestCaseExpanded = (test_case_name: string) => expandedTestCaseNames.includes(test_case_name);

  const expandPre = (e) => {
    if(e.currentTarget.classList.contains('expandedPre')){
      e.currentTarget.classList.remove("expandedPre")
      e.currentTarget.classList.add('collapsedPre')

    } else {
      e.currentTarget.classList.remove('collapsedPre')
      e.currentTarget.classList.add('expandedPre')
    }
  }

  const [activeSortIndex, setActiveSortIndex] = React.useState<number>(0);
  const [activeSortDirection, setActiveSortDirection] = React.useState<'asc' | 'desc' >('desc');
  
  const getSortParams = (columnIndex: number): ThProps['sort'] => ({
    sortBy: {
      index: activeSortIndex,
      direction: activeSortDirection,
      defaultDirection: 'asc' // starting sort direction when first sorting a column. Defaults to 'asc'
    },
    onSort: (_event, index, direction) => {
      setActiveSortIndex(index);
      setActiveSortDirection(direction);
    },
    columnIndex
  });

  const getSortableRowValues = (flakey: Flakey): (string | number)[] => {
    const {  average_impact, suite_name } = flakey;
    return [ average_impact, suite_name];
  };

  const onSortFn = (a:Flakey, b:Flakey):number => {
    
    const aValue = getSortableRowValues(a)[activeSortIndex];
    const bValue = getSortableRowValues(b)[activeSortIndex];
    if (typeof aValue === 'number') {
      // Numeric sort
      if (activeSortDirection === 'asc') {
        return (aValue as number) - (bValue as number);
      }
      return (bValue as number) - (aValue as number);
    } else {
      // String sort
      if (activeSortDirection === 'asc') {
        return (aValue as string).localeCompare(bValue as string);
      }
      return (bValue as string).localeCompare(aValue as string);
    }
  }

  const countFailingSuites = (suite:Flakey):number => {
    const sum = suite.test_cases.reduce((accumulator, object) => {
      return accumulator + object.count;
    }, 0);
    return sum
  }

  return (
    <InnerScrollContainer>
      <TableComposable aria-label="Nested column headers with expandable rows table" gridBreakPoint="">
        <Thead hasNestedHeader>
          <Tr>
            <Th rowSpan={2} />
            <Th width={50} >
              {columnNames.suite_name}
            </Th>
            <Th width={10} sort={getSortParams(0)}>
              {columnNames.average_impact}
            </Th>
          </Tr>
        </Thead>
        {teams.sort(onSortFn).map((suite, rowIndex) => (
          <Tbody key={suite.suite_name+rowIndex} isExpanded={isSuiteExpanded(suite.suite_name)}>
            <Tr>
              <Td expand={{ rowIndex, isExpanded: isSuiteExpanded(suite.suite_name), onToggle: () => setSuiteExpanded(suite.suite_name, !isSuiteExpanded(suite.suite_name))}}/>
              <Td dataLabel={columnNames.name}>{suite.suite_name}</Td>
              <Td dataLabel={columnNames.count}>{suite.average_impact.toFixed(2)}%</Td>
            </Tr>
            <Tr isExpanded={isSuiteExpanded(suite.suite_name)} className='pf-px-xl'>
              <Td colSpan={3}>
                <ExpandableRowContent>
                  <Grid hasGutter>
                    <GridItem span={3}>
                      <Card className='card-no-border'>
                        <CardTitle component="h3" style={{color: "red"}}>Overall Impact</CardTitle>
                        <CardBody>{suite.average_impact.toFixed(2)}%</CardBody>
                      </Card>
                      <Card className='card-no-border'>
                        <CardTitle component="h3">Total count of failed test cases</CardTitle>
                        <CardBody>{isSuiteExpanded(suite.suite_name) ? countFailingSuites(suite): ""}</CardBody>
                      </Card>
                    </GridItem>
                    <GridItem span={9}>
                      <Card className='card-no-border'>
                        <CardTitle>Failing test cases</CardTitle>
                        <TableComposable aria-label="Error messages" variant="compact">
                          <Thead>
                            <Tr>
                              <Th width={80} rowSpan={2} />
                              <Th>Test Name</Th>
                              <Th width={10}>
                                Impact
                              </Th>
                              <Th width={10}>
                                Count
                              </Th>
                            </Tr>
                          </Thead>
                          {suite.test_cases && suite.test_cases.map((test_case, tc_idx) => (
                            <Tbody key={test_case.name+tc_idx}>
                              <Tr>
                                <Td expand={{ rowIndex, isExpanded: isTestCaseExpanded(test_case.name), onToggle: () => setTestCaseExpanded(test_case.name, !isTestCaseExpanded(test_case.name))}}/>
                                <Td>
                                  {test_case.name}
                                </Td>
                                <Td>
                                  {test_case.test_case_impact}
                                </Td>
                                <Td>
                                  {test_case.count}
                                </Td>
                              </Tr>
                              <Tr isExpanded={isTestCaseExpanded(test_case.name)}>
                                <Td></Td>
                                <Td colSpan={3}>
                                  <ExpandableRowContent>

                                    <TableComposable aria-label="Error messages" variant="compact">
                                      <Thead>
                                        <Tr>
                                          <Th />
                                          <Th width={10}>Job URL</Th>
                                          <Th width={70}>Error</Th>
                                          <Th width={20}>Failure Dates</Th>
                                        </Tr>
                                      </Thead>
                                      {test_case.messages && test_case.messages.map((message, m_idx) => (
                                        <Tbody key={message.job_id+m_idx}>
                                        <Tr>
                                          <Td></Td>
                                          <Td dataLabel={columnNames.job_id}>
                                            <a href={message.job_url} rel="noreferrer noopener" target='_blank'>{message.job_id}</a>
                                          </Td>
                                          <Td dataLabel="Error messages" onClick={expandPre} className='collapsedPre'>
                                            <p  style={{textAlign: 'center', color: 'var(--pf-global--link--Color)', cursor: "pointer"}}><u>Show error</u></p>
                                            <pre>
                                              {message.error_message}
                                            </pre>
                                          </Td>
                                          <Td dataLabel={columnNames.failure_date}>
                                            {message.failure_date}
                                          </Td>
                                        </Tr>
                                        </Tbody>
                                        ))
                                      }
                                    </TableComposable>
                                  </ExpandableRowContent>

                                </Td>
                              </Tr>
                            </Tbody>
                          ))}
                        </TableComposable>
                      </Card>
                    </GridItem>
                  </Grid>
                </ExpandableRowContent>
              </Td>
            </Tr>
          </Tbody>
        ))}
      </TableComposable>  
    </InnerScrollContainer>
  );
};

const FlakeyTests: React.FunctionComponent = () => {
  const [toggles, setToggles] = React.useState<any>([])
  const [selectedSuite, setSelectedSuite] = React.useState<string>('')
  const [data, setData] = React.useState<any>([])
  const [pieData, setPieData] = React.useState<any>([])
  const [barData, setBarData] = React.useState<any>([])
  const [selectedRepo, setSelectedRepo] = React.useState<string>('')
  const [repositories, setRepositories] = React.useState<Array<any>>([])
  const [jobs, setJobs] = React.useState<Array<any>>([])
  const [selectedJob, setSelectedJob] = React.useState<string>('')
  const [startDate, setStartDate] = React.useState<Date | undefined>(undefined)
  const [endDate, setEndDate] = React.useState<Date | undefined>(undefined)
  const currentTeam = useSelector((state: any) => state.teams.Team);

  const countSuiteFailures = (suites) => {
    return suites.map((suite) => {
      const c = suite.test_cases.reduce(function (acc, obj) { return acc + obj.count; }, 0);
      return {suite_name: suite.suite_name, count: c}
    })
  }

  React.useEffect(() => {
    if(!currentTeam){ return }
    getRepositoriesWithJobs(currentTeam).then( res => {
      if(res.code == 200){
        setRepositories(res.data)
      } else {
        console.error("Cannot get repositories and jobs")
      }
    })
  }, [currentTeam]);

  React.useEffect(() => {
    if(data){
      const organizedData = countSuiteFailures(data)
      setToggles(organizedData)
      setPieData(organizedData)
    }

  }, [data]);

  const onSuiteSelect = (value) => {
    setSelectedSuite(value)
  }

  const onDataFilter = (suite:Flakey) => {
    return suite.suite_name == selectedSuite || selectedSuite == '' || selectedSuite == 'All failures'
  }

  const onRepoSelect = (value) => {
    setSelectedRepo(value)
    setSelectedJob('')
    const j = repositories.filter(repo => repo.Repository.Name == value )
    if(j && j[0] && j[0].jobs){
      setJobs(j[0].jobs)
    }
  }

  const onJobSelect = (value) => {
    if(value != selectedJob){
      setSelectedJob(value)
    }
  }

  const onDatesChange = (value:string, date:Date, name:string) => {
    if(name=='start-date'){ 
      console.log(value, date.toISOString())
      setStartDate(date)
    }
    if(name=='end-date'){ 
      console.log(value, date.toISOString())
      setEndDate(date) 
    }
  }

  const fetchData = () => {
    if(startDate && endDate && selectedRepo && selectedRepo){
      getFlakyData(currentTeam, selectedJob, selectedRepo, startDate.toISOString(), endDate.toISOString(), "redhat-appstudio").then(res => {
        if(res.code == 200){
          const impact:FlakeyObject = res.data
          if(impact && impact.suites){
            setData(res.data.suites)
          }
        } else {
          console.log("error", res)
        }
      })
      getGlobalImpactData(currentTeam, selectedJob, selectedRepo, startDate.toISOString(), endDate.toISOString(), "redhat-appstudio").then(res => {
        if(res.code == 200){
          const impact:any = res.data
          console.log(impact)
          if(impact && impact.length>0){
            setBarData(impact.map(impact => { impact.Date = impact.Date.split(' ')[0]; return impact;}))
          }
        } else {
          console.log("error", res)
        }
      })
    }
  }

  return (
    <React.Fragment>
      <Header info="Observe the impact of the flaky tests that are affecting CI."></Header>
      <PageSection variant={PageSectionVariants.light}>
        <Title headingLevel="h3" size={TitleSizes['2xl']}>
          Flaky tests impacting CI
        </Title>
      </PageSection>
      <PageSection style={{minHeight: '120vh'}} variant={PageSectionVariants.default}>
        <div>
          <Grid hasGutter>
            <GridItem span={12}>
              <TextContent className='bg-white'>
                <Text component={TextVariants.h1}>Flaky tetst</Text>
              </TextContent>
            </GridItem>
            <GridItem span={12}>
              <Toolbar id="toolbar-items">
                <ToolbarContent>
                  <DropdownBasic selected={selectedRepo} toggles={
                    repositories.map( (repo, idx) => <DropdownItem key={idx} name={repo.Repository.Name}> {repo.Repository.Name} </DropdownItem> )
                  } onSelect={onRepoSelect} placeholder="Select a repository"></DropdownBasic>
                  <DropdownBasic selected={selectedJob} toggles={
                    jobs.map( (job, idx) => <DropdownItem key={idx} name={job}> {job} </DropdownItem> )
                  } onSelect={onJobSelect} placeholder="Select a repository"></DropdownBasic>
                  <DatePickerMinMax onChange={onDatesChange} selectedDate={startDate} name="start-date" ></DatePickerMinMax>
                  <DatePickerMinMax onChange={onDatesChange} selectedDate={endDate} name="end-date" ></DatePickerMinMax>
                  <Button variant="primary" onClick={fetchData}>Primary</Button>
                </ToolbarContent>
              </Toolbar>
            </GridItem>
            <GridItem>
              <Grid hasGutter className='bg-white'>
                <GridItem span={5}>
                  <div>
                    <Title headingLevel="h3">Count of failed tests by suite</Title>
                    <PieChart data={pieData} x="suite_name" y="count"></PieChart>
                  </div>
                </GridItem>
                <GridItem span={7}>
                  <div>
                    <Title headingLevel="h3">Impact on CI suite</Title>
                    <ImpactChart data={barData} x="Date" y="global_impact"></ImpactChart>
                  </div>
                </GridItem>
              </Grid>  
            </GridItem>
            <GridItem style={{clear: 'both', minHeight: '1em'}} span={12}>
            </GridItem>
            <GridItem span={12}>
              <Toolbar id="toolbar-items">
                <ToolbarContent>
                  <DropdownBasic selected={selectedSuite} toggles={
                    [ <DropdownItem key={'all'} name="All failures"> All failures</DropdownItem>,
                      ...toggles.sort((a,b) => b.count - a.count).map( (toggle, idx) => <DropdownItem key={idx} name={toggle.suite_name}> {toggle.suite_name} (<strong style={{color: 'red'}}>{toggle.count}</strong>) </DropdownItem> )
                    ]
                  } onSelect={onSuiteSelect} placeholder="Select a failing suite"></DropdownBasic>
                </ToolbarContent>
              </Toolbar>
            </GridItem>
            <GridItem span={12}>
              <ComposableTableNestedExpandable teams={data.filter(onDataFilter)}></ComposableTableNestedExpandable>
            </GridItem>
          </Grid>
        </div>
      </PageSection>
    </React.Fragment>
  )
}

export { FlakeyTests };
