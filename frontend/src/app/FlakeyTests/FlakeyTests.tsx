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
import { Stack, StackItem } from '@patternfly/react-core';
import { PageSection, PageSectionVariants, Title, TitleSizes } from '@patternfly/react-core';
import { Dropdown, DropdownToggle, DropdownItem } from '@patternfly/react-core';
import { ChartPie } from '@patternfly/react-charts';
import { Chart, ChartAxis, ChartLine, ChartGroup, ChartVoronoiContainer } from '@patternfly/react-charts';
import { Grid, GridItem, Flex, FlexItem } from '@patternfly/react-core';
import { Card, CardTitle, CardBody, CardFooter } from '@patternfly/react-core';
import { Toolbar, ToolbarItem, ToolbarContent } from '@patternfly/react-core';

const ImpactChart:React.FunctionComponent<{data, x, y}> = ({data, x, y}) => {
  const ref = useRef<HTMLDivElement>(null);
  const [width, setWidth] = useState(0);
  const [height, setHeight] = useState(0);

  useLayoutEffect(() => {
    if (ref.current) {
      console.log(ref.current.offsetWidth, ref.current.offsetHeight)
      setWidth(ref.current.offsetWidth * 0.8 - 10);
      setHeight(ref.current.offsetWidth * 0.8 * 0.4 -20);
    }
  }, []);

  return (
    <div style={{  width: '100%', height: '100%', boxShadow: "none" }} className={"pf-c-card"} ref={ref}>
      <div style={{ height: height + 'px', width: width + 'px', background: "white", boxShadow: "none" }}>
        <Chart
          ariaDesc="Average number of pets"
          ariaTitle="Bar chart example"
          containerComponent={<ChartVoronoiContainer labels={({ datum }) => `${datum.name}: ${datum.y}`} constrainToVisibleArea />}
          domain={{y: [0,50]}}
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
          <ChartGroup offset={11}>
            <ChartLine data={  data.map( (datum) => { return {"name": datum[x], "x": datum[x], "y": parseFloat(datum[y])}  }) }/>
         </ChartGroup>
        </Chart>
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
          style={{labels: { fontSize: 9}}}
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
      </div>
    </div>
    )
}

export const DropdownBasic: React.FunctionComponent<{toggles, onSelect, selected}> = ({toggles, onSelect, selected}) => {
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
      toggle={<DropdownToggle id="toggle-basic" onToggle={onToggle}>{selected == '' ? "Select a failing suite" : selected}</DropdownToggle>}
      isOpen={isOpen}
      dropdownItems={
        [ <DropdownItem key={'all'} name="All failures"> All failures</DropdownItem>,
          ...toggles.sort((a,b) => b.count - a.count).map( (toggle, idx) => <DropdownItem key={idx} name={toggle.suite_name}> {toggle.suite_name} (<strong style={{color: 'red'}}>{toggle.count}</strong>) </DropdownItem> )
        ]
      }
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

  const mockData: FlakeyObject = {
    "global_impact": 40.54054054054054,
    "git_organization": "redhat-appstudio",
    "repository_name": "infra-deployments",
    "job_name": "pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests",
    "suites": [
        {
            "suite_name": "build-service-suite Build service E2E tests",
            "status": "failed",
            "test_cases": [
                {
                    "name": "[It] [build-service-suite Build service E2E tests] test PaC component build when a new component without specified branch is created and with visibility private a related PipelineRun should be deleted after deleting the component [build, HACBS, github-webhook, pac-build, pipeline, image-controller, pac-custom-default-branch]",
                    "test_case_impact": 5.405405405405405,
                    "count": 2,
                    "messages": [
                        {
                            "job_id": "1724346055316738048",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2719/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724346055316738048",
                            "error_message": "Expected success, but got an error:\n    <context.deadlineExceededError>: \n    context deadline exceeded\n    {}",
                            "failure_date": "2023-11-14T08:38:31Z"
                        },
                        {
                            "job_id": "1724016536861020160",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2738/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724016536861020160",
                            "error_message": "Expected success, but got an error:\n    <context.deadlineExceededError>: \n    context deadline exceeded\n    {}",
                            "failure_date": "2023-11-13T10:49:07Z"
                        }
                    ]
                },
                {
                    "name": "[It] [build-service-suite Build service E2E tests] PLNSRVCE-799 - test pipeline selector default Pipeline bundle should be used and no additional Pipeline params should be added to the PipelineRun if one of the WhenConditions does not match [build, HACBS, pipeline-selector]",
                    "test_case_impact": 5.405405405405405,
                    "count": 2,
                    "messages": [
                        {
                            "job_id": "1724095820627709952",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2739/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724095820627709952",
                            "error_message": "Expected success, but got an error:\n    <context.deadlineExceededError>: \n    context deadline exceeded\n    {}",
                            "failure_date": "2023-11-13T16:04:10Z"
                        },
                        {
                            "job_id": "1724069906363715584",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2629/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724069906363715584",
                            "error_message": "Expected success, but got an error:\n    <context.deadlineExceededError>: \n    context deadline exceeded\n    {}",
                            "failure_date": "2023-11-13T14:21:11Z"
                        }
                    ]
                },
                {
                    "name": "[It] [build-service-suite Build service E2E tests] test pac with multiple components using same repository when components are created in same namespace the PipelineRun should eventually finish successfully for component go-component-wcve [build, HACBS, pac-build, multi-component]",
                    "test_case_impact": 2.7027027027027026,
                    "count": 1,
                    "messages": [
                        {
                            "job_id": "1724390754152878080",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2741/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724390754152878080",
                            "error_message": "Expected success, but got an error:\n    <*errors.errorString | 0xc00109f270>: \n    \n    init container prepare: \n    2023/11/14 11:53:01 Entrypoint initialization\n    \n    init container place-scripts: \n    2023/11/14 11:53:02 Decoded script /tekton/scripts/script-0-bbf9w\n    2023/11/14 11:53:02 Decoded script /tekton/scripts/script-1-kwnxz\n    \n    container step-clone: \n    + '[' true = true ']'\n    + '[' -f /workspace/basic-auth/.git-credentials ']'\n    + '[' -f /workspace/basic-auth/.gitconfig ']'\n    + cp /workspace/basic-auth/.git-credentials /tekton/home/.git-credentials\n    + cp /workspace/basic-auth/.gitconfig /tekton/home/.gitconfig\n    + chmod 400 /tekton/home/.git-credentials\n    + chmod 400 /tekton/home/.gitconfig\n    + '[' false = true ']'\n    + CHECKOUT_DIR=/workspace/output/source\n    + '[' true = true ']'\n    + cleandir\n    + '[' -d /workspace/output/source ']'\n    + test -z ''\n    + test -z ''\n    + test -z ''\n    + /ko-app/git-init -url=https://github.com/redhat-appstudio-qe/sample-multi-component -revision=05ab65ffc94b843e8543556c20bd16647ffc3220 -refspec= -path=/workspace/output/source -sslVerify=true -submodules=true -depth=1 -sparseCheckoutDirectories=\n    {\"level\":\"info\",\"ts\":1699962800.4609172,\"caller\":\"git/git.go:178\",\"msg\":\"Successfully cloned https://github.com/redhat-appstudio-qe/sample-multi-component @ 05ab65ffc94b843e8543556c20bd16647ffc3220 (grafted, HEAD) in path /workspace/output/source\"}\n    {\"level\":\"info\",\"ts\":1699962800.491771,\"caller\":\"git/git.go:217\",\"msg\":\"Successfully initialized and updated submodules in path /workspace/output/source\"}\n    + cd /workspace/output/source\n    ++ git rev-parse HEAD\n    + RESULT_SHA=05ab65ffc94b843e8543556c20bd16647ffc3220\n    + EXIT_CODE=0\n    + '[' 0 '!=' 0 ']'\n    + printf %!s(MISSING) 05ab65ffc94b843e8543556c20bd16647ffc3220\n    + printf %!s(MISSING) https://github.com/redhat-appstudio-qe/sample-multi-component\n    + '[' false = true ']'\n    \n    container step-symlink-check: \n    Running symlink check\n    \n    init container prepare: \n    2023/11/14 11:52:40 Entrypoint initialization\n    \n    init container place-scripts: \n    2023/11/14 11:52:41 Decoded script /tekton/scripts/script-0-hgqxq\n    \n    container step-init: \n    Build Initialize: quay.io/redhat-appstudio-qe/build-e2e-cpme-tenant/build-suite-positive-mc-nupu/go-component-wcve:on-pr-05ab65ffc94b843e8543556c20bd16647ffc3220\n    \n    Determine if Image Already Exists\n    \n    init container prepare: \n    2023/11/14 11:53:26 Entrypoint initialization\n    \n    init container place-scripts: \n    2023/11/14 11:53:27 Decoded script /tekton/scripts/script-0-25ssv\n    \n    init container working-dir-initializer: \n    \n    container step-sast-snyk-check: \n    {\"result\":\"SKIPPED\",\"timestamp\":\"1699962832\",\"note\":\"Task sast-snyk-check skipped: If you wish to use the Snyk code SAST task, please create a secret named snyk-secret with the key snyk_token containing the Snyk token.\",\"namespace\":\"default\",\"successes\":0,\"failures\":0,\"warnings\":0}\n    \n    init container prepare: \n    2023/11/14 11:53:56 Entrypoint initialization\n    \n    init container place-scripts: \n    2023/11/14 11:53:56 Decoded script /tekton/scripts/script-0-rlvw2\n    \n    container step-appstudio-summary: \n    \n    Build Summary:\n    \n    Build repository: https://github.com/redhat-appstudio-qe/sample-multi-component?rev=05ab65ffc94b843e8543556c20bd16647ffc3220\n    \n    End Summary\n    \n    {\n        s: \"\\ninit container prepare: \\n2023/11/14 11:53:01 Entrypoint initialization\\n\\ninit container place-scripts: \\n2023/11/14 11:53:02 Decoded script /tekton/scripts/script-0-bbf9w\\n2023/11/14 11:53:02 Decoded script /tekton/scripts/script-1-kwnxz\\n\\ncontainer step-clone: \\n+ '[' true = true ']'\\n+ '[' -f /workspace/basic-auth/.git-credentials ']'\\n+ '[' -f /workspace/basic-auth/.gitconfig ']'\\n+ cp /workspace/basic-auth/.git-credentials /tekton/home/.git-credentials\\n+ cp /workspace/basic-auth/.gitconfig /tekton/home/.gitconfig\\n+ chmod 400 /tekton/home/.git-credentials\\n+ chmod 400 /tekton/home/.gitconfig\\n+ '[' false = true ']'\\n+ CHECKOUT_DIR=/workspace/output/source\\n+ '[' true = true ']'\\n+ cleandir\\n+ '[' -d /workspace/output/source ']'\\n+ test -z ''\\n+ test -z ''\\n+ test -z ''\\n+ /ko-app/git-init -url=https://github.com/redhat-appstudio-qe/sample-multi-component -revision=05ab65ffc94b843e8543556c20bd16647ffc3220 -refspec= -path=/workspace/output/source -sslVerify=true -submodules=true -depth=1 -sparseCheckoutDirectories=\\n{\\\"level\\\":\\\"info\\\",\\\"ts\\\":1699962800.4609172,\\\"caller\\\":\\\"git/git.go:178\\\",\\\"msg\\\":\\\"Successfully cloned https://github.com/redhat-appstudio-qe/sample-multi-component @ 05ab65ffc94b843e8543556c20bd16647ffc3220 (grafted, HEAD) in path /workspace/output/source\\\"}\\n{\\\"level\\\":\\\"info\\\",\\\"ts\\\":1699962800.491771,\\\"caller\\\":\\\"git/git.go:217\\\",\\\"msg\\\":\\\"Successfully initialized and updated submodules in path /workspace/output/source\\\"}\\n+ cd /workspace/output/source\\n++ git rev-parse HEAD\\n+ RESULT_SHA=05ab65ffc94b843e8543556c20bd16647ffc3220\\n+ EXIT_CODE=0\\n+ '[' 0 '!=' 0 ']'\\n+ printf %!s(MISSING) 05ab65ffc94b843e8543556c20bd16647ffc3220\\n+ printf %!s(MISSING) https://github.com/redhat-appstudio-qe/sample-multi-component\\n+ '[' false = true ']'\\n\\ncontainer step-symlink-check: \\nRunning symlink check\\n\\ninit container prepare: \\n2023/11/14 11:52:40 Entrypoint initialization\\n\\ninit container place-scripts: \\n2023/11/14 11:52:41 Decoded script /tekton/scripts/script-0-hgqxq\\n\\ncontainer step-init: \\nBuild Initialize: quay.io/redhat-appstudio-qe/build-e2e-cpme-tenant/build-suite-positive-mc-nupu/go-component-wcve:on-pr-05ab65ffc94b843e8543556c20bd16647ffc3220\\n\\nDetermine if Image Already Exists\\n\\ninit container prepare: \\n2023/11/14 11:53:26 Entrypoint initialization\\n\\ninit container place-scripts: \\n2023/11/14 11:53:27 Decoded script /tekton/scripts/script-0-25ssv\\n\\ninit container working-dir-initializer: \\n\\ncontainer step-sast-snyk-check: \\n{\\\"result\\\":\\\"SKIPPED\\\",\\\"timestamp\\\":\\\"1699962832\\\",\\\"note\\\":\\\"Task sast-snyk-check skipped: If you wish to use the Snyk code SAST task, please create a secret named snyk-secret with the key snyk_token containing the Snyk token.\\\",\\\"namespace\\\":\\\"default\\\",\\\"successes\\\":0,\\\"failures\\\":0,\\\"warnings\\\":0}\\n\\ninit container prepare: \\n2023/11/14 11:53:56 Entrypoint initialization\\n\\ninit container place-scripts: \\n2023/11/14 11:53:56 Decoded script /tekton/scripts/script-0-rlvw2\\n\\ncontainer step-appstudio-summary: \\n\\nBuild Summary:\\n\\nBuild repository: https://github.com/redhat-appstudio-qe/sample-multi-component?rev=05ab65ffc94b843e8543556c20bd16647ffc3220\\n\\nEnd Summary\\n\",\n    }",
                            "failure_date": "2023-11-14T11:36:08Z"
                        }
                    ]
                },
                {
                    "name": "[It] [build-service-suite Build service E2E tests] test PaC component build when the component is removed and recreated (with the same name in the same namespace) should no longer lead to a creation of a PaC PR [build, HACBS, github-webhook, pac-build, pipeline, image-controller]",
                    "test_case_impact": 2.7027027027027026,
                    "count": 1,
                    "messages": [
                        {
                            "job_id": "1724398958685458432",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2698/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724398958685458432",
                            "error_message": "Expected success, but got an error:\n    <context.deadlineExceededError>: \n    context deadline exceeded\n    {}",
                            "failure_date": "2023-11-14T12:08:44Z"
                        }
                    ]
                }
            ],
            "average_impact": 16.216216216216218
        },
        {
            "suite_name": "build-service-suite Build templates E2E test",
            "status": "failed",
            "test_cases": [
                {
                    "name": "[It] [build-service-suite Build templates E2E test] HACBS pipelines should eventually finish successfully for component with Git source URL https://github.com/redhat-appstudio-qe/devfile-sample-python-basic.git [build, HACBS, pipeline, build-templates-e2e]",
                    "test_case_impact": 2.7027027027027026,
                    "count": 1,
                    "messages": [
                        {
                            "job_id": "1724386610788700160",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2719/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724386610788700160",
                            "error_message": "Expected success, but got an error:\n    <*errors.errorString | 0xc001669870>: \n    \n    init container prepare: \n    2023/11/14 11:35:11 Entrypoint initialization\n    \n    init container place-scripts: \n    2023/11/14 11:35:13 Decoded script /tekton/scripts/script-0-vs9f6\n    2023/11/14 11:35:13 Decoded script /tekton/scripts/script-1-q7jrv\n    \n    container step-clone: \n    + '[' false = true ']'\n    + '[' false = true ']'\n    + CHECKOUT_DIR=/workspace/output/source\n    + '[' true = true ']'\n    + cleandir\n    + '[' -d /workspace/output/source ']'\n    + test -z ''\n    + test -z ''\n    + test -z ''\n    + /ko-app/git-init -url=https://github.com/redhat-appstudio-qe/devfile-sample-python-basic.git -revision=main -refspec= -path=/workspace/output/source -sslVerify=true -submodules=true -depth=1 -sparseCheckoutDirectories=\n    {\"level\":\"info\",\"ts\":1699961719.3569608,\"caller\":\"git/git.go:178\",\"msg\":\"Successfully cloned https://github.com/redhat-appstudio-qe/devfile-sample-python-basic.git @ 7c630e200f40ba457ed508e7f6401d03fb50982d (grafted, HEAD, origin/main) in path /workspace/output/source\"}\n    {\"level\":\"info\",\"ts\":1699961719.3835516,\"caller\":\"git/git.go:217\",\"msg\":\"Successfully initialized and updated submodules in path /workspace/output/source\"}\n    + cd /workspace/output/source\n    ++ git rev-parse HEAD\n    + RESULT_SHA=7c630e200f40ba457ed508e7f6401d03fb50982d\n    + EXIT_CODE=0\n    + '[' 0 '!=' 0 ']'\n    + printf %!s(MISSING) 7c630e200f40ba457ed508e7f6401d03fb50982d\n    + printf %!s(MISSING) https://github.com/redhat-appstudio-qe/devfile-sample-python-basic.git\n    + '[' false = true ']'\n    \n    container step-symlink-check: \n    Running symlink check\n    \n    init container prepare: \n    2023/11/14 11:34:49 Entrypoint initialization\n    \n    init container place-scripts: \n    2023/11/14 11:34:50 Decoded script /tekton/scripts/script-0-m88s5\n    \n    container step-init: \n    Build Initialize: quay.io/redhat-appstudio-qe/build-e2e-cdlg-tenant/test-app-qsly/devfile-sample-python-basic-xxqf:build-b34d2-1699961683\n    \n    Determine if Image Already Exists\n    \n    init container prepare: \n    2023/11/14 11:35:24 Entrypoint initialization\n    \n    init container place-scripts: \n    2023/11/14 11:35:25 Decoded script /tekton/scripts/script-0-dwxnq\n    \n    init container working-dir-initializer: \n    \n    container step-sast-snyk-check: \n    {\"result\":\"SKIPPED\",\"timestamp\":\"1699961745\",\"note\":\"Task sast-snyk-check skipped: If you wish to use the Snyk code SAST task, please create a secret named snyk-secret with the key snyk_token containing the Snyk token.\",\"namespace\":\"default\",\"successes\":0,\"failures\":0,\"warnings\":0}\n    \n    init container prepare: \n    2023/11/14 11:35:54 Entrypoint initialization\n    \n    init container place-scripts: \n    2023/11/14 11:35:55 Decoded script /tekton/scripts/script-0-6r5lk\n    \n    container step-appstudio-summary: \n    \n    Build Summary:\n    \n    Build repository: https://github.com/redhat-appstudio-qe/devfile-sample-python-basic.git?rev=7c630e200f40ba457ed508e7f6401d03fb50982d\n    \n    End Summary\n    \n    {\n        s: \"\\ninit container prepare: \\n2023/11/14 11:35:11 Entrypoint initialization\\n\\ninit container place-scripts: \\n2023/11/14 11:35:13 Decoded script /tekton/scripts/script-0-vs9f6\\n2023/11/14 11:35:13 Decoded script /tekton/scripts/script-1-q7jrv\\n\\ncontainer step-clone: \\n+ '[' false = true ']'\\n+ '[' false = true ']'\\n+ CHECKOUT_DIR=/workspace/output/source\\n+ '[' true = true ']'\\n+ cleandir\\n+ '[' -d /workspace/output/source ']'\\n+ test -z ''\\n+ test -z ''\\n+ test -z ''\\n+ /ko-app/git-init -url=https://github.com/redhat-appstudio-qe/devfile-sample-python-basic.git -revision=main -refspec= -path=/workspace/output/source -sslVerify=true -submodules=true -depth=1 -sparseCheckoutDirectories=\\n{\\\"level\\\":\\\"info\\\",\\\"ts\\\":1699961719.3569608,\\\"caller\\\":\\\"git/git.go:178\\\",\\\"msg\\\":\\\"Successfully cloned https://github.com/redhat-appstudio-qe/devfile-sample-python-basic.git @ 7c630e200f40ba457ed508e7f6401d03fb50982d (grafted, HEAD, origin/main) in path /workspace/output/source\\\"}\\n{\\\"level\\\":\\\"info\\\",\\\"ts\\\":1699961719.3835516,\\\"caller\\\":\\\"git/git.go:217\\\",\\\"msg\\\":\\\"Successfully initialized and updated submodules in path /workspace/output/source\\\"}\\n+ cd /workspace/output/source\\n++ git rev-parse HEAD\\n+ RESULT_SHA=7c630e200f40ba457ed508e7f6401d03fb50982d\\n+ EXIT_CODE=0\\n+ '[' 0 '!=' 0 ']'\\n+ printf %!s(MISSING) 7c630e200f40ba457ed508e7f6401d03fb50982d\\n+ printf %!s(MISSING) https://github.com/redhat-appstudio-qe/devfile-sample-python-basic.git\\n+ '[' false = true ']'\\n\\ncontainer step-symlink-check: \\nRunning symlink check\\n\\ninit container prepare: \\n2023/11/14 11:34:49 Entrypoint initialization\\n\\ninit container place-scripts: \\n2023/11/14 11:34:50 Decoded script /tekton/scripts/script-0-m88s5\\n\\ncontainer step-init: \\nBuild Initialize: quay.io/redhat-appstudio-qe/build-e2e-cdlg-tenant/test-app-qsly/devfile-sample-python-basic-xxqf:build-b34d2-1699961683\\n\\nDetermine if Image Already Exists\\n\\ninit container prepare: \\n2023/11/14 11:35:24 Entrypoint initialization\\n\\ninit container place-scripts: \\n2023/11/14 11:35:25 Decoded script /tekton/scripts/script-0-dwxnq\\n\\ninit container working-dir-initializer: \\n\\ncontainer step-sast-snyk-check: \\n{\\\"result\\\":\\\"SKIPPED\\\",\\\"timestamp\\\":\\\"1699961745\\\",\\\"note\\\":\\\"Task sast-snyk-check skipped: If you wish to use the Snyk code SAST task, please create a secret named snyk-secret with the key snyk_token containing the Snyk token.\\\",\\\"namespace\\\":\\\"default\\\",\\\"successes\\\":0,\\\"failures\\\":0,\\\"warnings\\\":0}\\n\\ninit container prepare: \\n2023/11/14 11:35:54 Entrypoint initialization\\n\\ninit container place-scripts: \\n2023/11/14 11:35:55 Decoded script /tekton/scripts/script-0-6r5lk\\n\\ncontainer step-appstudio-summary: \\n\\nBuild Summary:\\n\\nBuild repository: https://github.com/redhat-appstudio-qe/devfile-sample-python-basic.git?rev=7c630e200f40ba457ed508e7f6401d03fb50982d\\n\\nEnd Summary\\n\",\n    }",
                            "failure_date": "2023-11-14T11:19:40Z"
                        }
                    ]
                },
                {
                    "name": "[It] [build-service-suite Build templates E2E test] HACBS pipelines when Pipeline Results are stored for component with Git source URL https://github.com/redhat-appstudio-qe/devfile-sample-python-basic.git should have Pipeline Logs [build, HACBS, pipeline]",
                    "test_case_impact": 2.7027027027027026,
                    "count": 1,
                    "messages": [
                        {
                            "job_id": "1724439144450494464",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2748/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724439144450494464",
                            "error_message": "Timed out after 120.000s.\ntimed out when getting logs for PipelineRun build-e2e-fual-tenant/devfile-sample-python-basic-gjmy-dwvwl\nExpected success, but got an error:\n    <*errors.errorString | 0xc001e82870>: \n    logs for PipelineRun build-e2e-fual-tenant/devfile-sample-python-basic-gjmy-dwvwl are empty\n    {\n        s: \"logs for PipelineRun build-e2e-fual-tenant/devfile-sample-python-basic-gjmy-dwvwl are empty\",\n    }",
                            "failure_date": "2023-11-14T14:48:25Z"
                        }
                    ]
                }
            ],
            "average_impact": 5.405405405405405
        },
        {
            "suite_name": "byoc-suite",
            "status": "failed",
            "test_cases": [
                {
                    "name": "[It] [byoc-suite] Deploy RHTAP sample application into a Kubernetes cluster provided by user waits component pipeline to finish [byoc]",
                    "test_case_impact": 2.7027027027027026,
                    "count": 1,
                    "messages": [
                        {
                            "job_id": "1724310537740750848",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2738/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724310537740750848",
                            "error_message": "Expected success, but got an error:\n    <*errors.errorString | 0xc0019b9c20>: \n    \n    init container prepare: \n    2023/11/14 06:36:57 Entrypoint initialization\n    \n    init container place-scripts: \n    2023/11/14 06:36:58 Decoded script /tekton/scripts/script-0-dpdtt\n    2023/11/14 06:36:58 Decoded script /tekton/scripts/script-1-xwcsr\n    2023/11/14 06:36:58 Decoded script /tekton/scripts/script-2-mbtpt\n    2023/11/14 06:36:58 Decoded script /tekton/scripts/script-3-7hwcv\n    2023/11/14 06:36:58 Decoded script /tekton/scripts/script-4-4rxk6\n    2023/11/14 06:36:58 Decoded script /tekton/scripts/script-5-t8kj4\n    2023/11/14 06:36:58 Decoded script /tekton/scripts/script-6-8qskp\n    \n    init container working-dir-initializer: \n    \n    container step-build: \n    \n    container step-sbom-syft-generate: \n    \n    container step-analyse-dependencies-java-sbom: \n    \n    container step-merge-syft-sboms: \n    \n    container step-merge-cachi2-sbom: \n    \n    container step-create-purl-sbom: \n    \n    container step-inject-sbom-and-push: \n    \n    container step-upload-sbom: \n    \n    init container prepare: \n    2023/11/14 06:36:38 Entrypoint initialization\n    \n    init container place-scripts: \n    2023/11/14 06:36:38 Decoded script /tekton/scripts/script-0-dtrhs\n    2023/11/14 06:36:38 Decoded script /tekton/scripts/script-1-dszs4\n    \n    container step-clone: \n    + '[' false = true ']'\n    + '[' false = true ']'\n    + CHECKOUT_DIR=/workspace/output/source\n    + '[' true = true ']'\n    + cleandir\n    + '[' -d /workspace/output/source ']'\n    + test -z ''\n    + test -z ''\n    + test -z ''\n    + /ko-app/git-init -url=https://github.com/devfile-samples/devfile-sample-code-with-quarkus -revision=main -refspec= -path=/workspace/output/source -sslVerify=true -submodules=true -depth=1 -sparseCheckoutDirectories=\n    {\"level\":\"info\",\"ts\":1699943812.3555655,\"caller\":\"git/git.go:178\",\"msg\":\"Successfully cloned https://github.com/devfile-samples/devfile-sample-code-with-quarkus @ 1815bc17f90687ff54b0d73a9486e10e4a69d5aa (grafted, HEAD, origin/main) in path /workspace/output/source\"}\n    {\"level\":\"info\",\"ts\":1699943812.383541,\"caller\":\"git/git.go:217\",\"msg\":\"Successfully initialized and updated submodules in path /workspace/output/source\"}\n    + cd /workspace/output/source\n    ++ git rev-parse HEAD\n    + RESULT_SHA=1815bc17f90687ff54b0d73a9486e10e4a69d5aa\n    + EXIT_CODE=0\n    + '[' 0 '!=' 0 ']'\n    + printf %!s(MISSING) 1815bc17f90687ff54b0d73a9486e10e4a69d5aa\n    + printf %!s(MISSING) https://github.com/devfile-samples/devfile-sample-code-with-quarkus\n    + '[' false = true ']'\n    \n    container step-symlink-check: \n    Running symlink check\n    \n    init container prepare: \n    2023/11/14 06:36:19 Entrypoint initialization\n    \n    init container place-scripts: \n    2023/11/14 06:36:19 Decoded script /tekton/scripts/script-0-7jzh9\n    \n    container step-init: \n    Build Initialize: quay.io/redhat-appstudio-qe/byoc-vetg-tenant/byoc-app-uasg/lbvm:build-0db21-1699943772\n    \n    Determine if Image Already Exists\n    \n    init container prepare: \n    2023/11/14 06:37:44 Entrypoint initialization\n    \n    init container place-scripts: \n    2023/11/14 06:37:46 Decoded script /tekton/scripts/script-0-8lcv6\n    \n    container step-appstudio-summary: \n    \n    Build Summary:\n    \n    Build repository: https://github.com/devfile-samples/devfile-sample-code-with-quarkus?rev=1815bc17f90687ff54b0d73a9486e10e4a69d5aa\n    \n    End Summary\n    \n    {\n        s: \"\\ninit container prepare: \\n2023/11/14 06:36:57 Entrypoint initialization\\n\\ninit container place-scripts: \\n2023/11/14 06:36:58 Decoded script /tekton/scripts/script-0-dpdtt\\n2023/11/14 06:36:58 Decoded script /tekton/scripts/script-1-xwcsr\\n2023/11/14 06:36:58 Decoded script /tekton/scripts/script-2-mbtpt\\n2023/11/14 06:36:58 Decoded script /tekton/scripts/script-3-7hwcv\\n2023/11/14 06:36:58 Decoded script /tekton/scripts/script-4-4rxk6\\n2023/11/14 06:36:58 Decoded script /tekton/scripts/script-5-t8kj4\\n2023/11/14 06:36:58 Decoded script /tekton/scripts/script-6-8qskp\\n\\ninit container working-dir-initializer: \\n\\ncontainer step-build: \\n\\ncontainer step-sbom-syft-generate: \\n\\ncontainer step-analyse-dependencies-java-sbom: \\n\\ncontainer step-merge-syft-sboms: \\n\\ncontainer step-merge-cachi2-sbom: \\n\\ncontainer step-create-purl-sbom: \\n\\ncontainer step-inject-sbom-and-push: \\n\\ncontainer step-upload-sbom: \\n\\ninit container prepare: \\n2023/11/14 06:36:38 Entrypoint initialization\\n\\ninit container place-scripts: \\n2023/11/14 06:36:38 Decoded script /tekton/scripts/script-0-dtrhs\\n2023/11/14 06:36:38 Decoded script /tekton/scripts/script-1-dszs4\\n\\ncontainer step-clone: \\n+ '[' false = true ']'\\n+ '[' false = true ']'\\n+ CHECKOUT_DIR=/workspace/output/source\\n+ '[' true = true ']'\\n+ cleandir\\n+ '[' -d /workspace/output/source ']'\\n+ test -z ''\\n+ test -z ''\\n+ test -z ''\\n+ /ko-app/git-init -url=https://github.com/devfile-samples/devfile-sample-code-with-quarkus -revision=main -refspec= -path=/workspace/output/source -sslVerify=true -submodules=true -depth=1 -sparseCheckoutDirectories=\\n{\\\"level\\\":\\\"info\\\",\\\"ts\\\":1699943812.3555655,\\\"caller\\\":\\\"git/git.go:178\\\",\\\"msg\\\":\\\"Successfully cloned https://github.com/devfile-samples/devfile-sample-code-with-quarkus @ 1815bc17f90687ff54b0d73a9486e10e4a69d5aa (grafted, HEAD, origin/main) in path /workspace/output/source\\\"}\\n{\\\"level\\\":\\\"info\\\",\\\"ts\\\":1699943812.383541,\\\"caller\\\":\\\"git/git.go:217\\\",\\\"msg\\\":\\\"Successfully initialized and updated submodules in path /workspace/output/source\\\"}\\n+ cd /workspace/output/source\\n++ git rev-parse HEAD\\n+ RESULT_SHA=1815bc17f90687ff54b0d73a9486e10e4a69d5aa\\n+ EXIT_CODE=0\\n+ '[' 0 '!=' 0 ']'\\n+ printf %!s(MISSING) 1815bc17f90687ff54b0d73a9486e10e4a69d5aa\\n+ printf %!s(MISSING) https://github.com/devfile-samples/devfile-sample-code-with-quarkus\\n+ '[' false = true ']'\\n\\ncontainer step-symlink-check: \\nRunning symlink check\\n\\ninit container prepare: \\n2023/11/14 06:36:19 Entrypoint initialization\\n\\ninit container place-scripts: \\n2023/11/14 06:36:19 Decoded script /tekton/scripts/script-0-7jzh9\\n\\ncontainer step-init: \\nBuild Initialize: quay.io/redhat-appstudio-qe/byoc-vetg-tenant/byoc-app-uasg/lbvm:build-0db21-1699943772\\n\\nDetermine if Image Already Exists\\n\\ninit container prepare: \\n2023/11/14 06:37:44 Entrypoint initialization\\n\\ninit container place-scripts: \\n2023/11/14 06:37:46 Decoded script /tekton/scripts/script-0-8lcv6\\n\\ncontainer step-appstudio-summary: \\n\\nBuild Summary:\\n\\nBuild repository: https://github.com/devfile-samples/devfile-sample-code-with-quarkus?rev=1815bc17f90687ff54b0d73a9486e10e4a69d5aa\\n\\nEnd Summary\\n\",\n    }",
                            "failure_date": "2023-11-14T06:17:22Z"
                        }
                    ]
                }
            ],
            "average_impact": 2.7027027027027026
        },
        {
            "suite_name": "integration-service-suite Integration Service E2E tests",
            "status": "failed",
            "test_cases": [
                {
                    "name": "[It] [integration-service-suite Integration Service E2E tests] with happy path for general flow of Integration service triggers a build PipelineRun [integration-service, HACBS]",
                    "test_case_impact": 2.7027027027027026,
                    "count": 1,
                    "messages": [
                        {
                            "job_id": "1724310537740750848",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2738/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724310537740750848",
                            "error_message": "Expected success, but got an error:\n    <*errors.errorString | 0xc0012b9c40>: \n    \n    init container prepare: \n    2023/11/14 06:36:10 Entrypoint initialization\n    \n    init container place-scripts: \n    2023/11/14 06:36:10 Decoded script /tekton/scripts/script-0-mwr7p\n    2023/11/14 06:36:10 Decoded script /tekton/scripts/script-1-vdnjx\n    \n    container step-clone: \n    + '[' false = true ']'\n    + '[' false = true ']'\n    + CHECKOUT_DIR=/workspace/output/source\n    + '[' true = true ']'\n    + cleandir\n    + '[' -d /workspace/output/source ']'\n    + test -z ''\n    + test -z ''\n    + test -z ''\n    + /ko-app/git-init -url=https://github.com/redhat-appstudio-qe/hacbs-test-project -revision=main -refspec= -path=/workspace/output/source -sslVerify=true -submodules=true -depth=1 -sparseCheckoutDirectories=\n    {\"level\":\"info\",\"ts\":1699943773.9714282,\"caller\":\"git/git.go:178\",\"msg\":\"Successfully cloned https://github.com/redhat-appstudio-qe/hacbs-test-project @ 34da5a8f51fba6a8b7ec75a727d3c72ebb5e1274 (grafted, HEAD, origin/main) in path /workspace/output/source\"}\n    {\"level\":\"info\",\"ts\":1699943774.0033467,\"caller\":\"git/git.go:217\",\"msg\":\"Successfully initialized and updated submodules in path /workspace/output/source\"}\n    + cd /workspace/output/source\n    ++ git rev-parse HEAD\n    + RESULT_SHA=34da5a8f51fba6a8b7ec75a727d3c72ebb5e1274\n    + EXIT_CODE=0\n    + '[' 0 '!=' 0 ']'\n    + printf %!s(MISSING) 34da5a8f51fba6a8b7ec75a727d3c72ebb5e1274\n    + printf %!s(MISSING) https://github.com/redhat-appstudio-qe/hacbs-test-project\n    + '[' false = true ']'\n    \n    container step-symlink-check: \n    Running symlink check\n    \n    init container prepare: \n    2023/11/14 06:35:50 Entrypoint initialization\n    \n    init container place-scripts: \n    2023/11/14 06:35:51 Decoded script /tekton/scripts/script-0-rldfw\n    \n    container step-init: \n    Build Initialize: quay.io/redhat-appstudio-qe/integration1-iegh-tenant/integ-app-olhh/hacbs-test-project-c1k1:build-abb81-1699943742\n    \n    Determine if Image Already Exists\n    \n    init container prepare: \n    2023/11/14 06:36:50 Entrypoint initialization\n    \n    init container place-scripts: \n    2023/11/14 06:36:51 Decoded script /tekton/scripts/script-0-m6zhw\n    \n    container step-appstudio-summary: \n    \n    Build Summary:\n    \n    Build repository: https://github.com/redhat-appstudio-qe/hacbs-test-project?rev=34da5a8f51fba6a8b7ec75a727d3c72ebb5e1274\n    \n    End Summary\n    \n    {\n        s: \"\\ninit container prepare: \\n2023/11/14 06:36:10 Entrypoint initialization\\n\\ninit container place-scripts: \\n2023/11/14 06:36:10 Decoded script /tekton/scripts/script-0-mwr7p\\n2023/11/14 06:36:10 Decoded script /tekton/scripts/script-1-vdnjx\\n\\ncontainer step-clone: \\n+ '[' false = true ']'\\n+ '[' false = true ']'\\n+ CHECKOUT_DIR=/workspace/output/source\\n+ '[' true = true ']'\\n+ cleandir\\n+ '[' -d /workspace/output/source ']'\\n+ test -z ''\\n+ test -z ''\\n+ test -z ''\\n+ /ko-app/git-init -url=https://github.com/redhat-appstudio-qe/hacbs-test-project -revision=main -refspec= -path=/workspace/output/source -sslVerify=true -submodules=true -depth=1 -sparseCheckoutDirectories=\\n{\\\"level\\\":\\\"info\\\",\\\"ts\\\":1699943773.9714282,\\\"caller\\\":\\\"git/git.go:178\\\",\\\"msg\\\":\\\"Successfully cloned https://github.com/redhat-appstudio-qe/hacbs-test-project @ 34da5a8f51fba6a8b7ec75a727d3c72ebb5e1274 (grafted, HEAD, origin/main) in path /workspace/output/source\\\"}\\n{\\\"level\\\":\\\"info\\\",\\\"ts\\\":1699943774.0033467,\\\"caller\\\":\\\"git/git.go:217\\\",\\\"msg\\\":\\\"Successfully initialized and updated submodules in path /workspace/output/source\\\"}\\n+ cd /workspace/output/source\\n++ git rev-parse HEAD\\n+ RESULT_SHA=34da5a8f51fba6a8b7ec75a727d3c72ebb5e1274\\n+ EXIT_CODE=0\\n+ '[' 0 '!=' 0 ']'\\n+ printf %!s(MISSING) 34da5a8f51fba6a8b7ec75a727d3c72ebb5e1274\\n+ printf %!s(MISSING) https://github.com/redhat-appstudio-qe/hacbs-test-project\\n+ '[' false = true ']'\\n\\ncontainer step-symlink-check: \\nRunning symlink check\\n\\ninit container prepare: \\n2023/11/14 06:35:50 Entrypoint initialization\\n\\ninit container place-scripts: \\n2023/11/14 06:35:51 Decoded script /tekton/scripts/script-0-rldfw\\n\\ncontainer step-init: \\nBuild Initialize: quay.io/redhat-appstudio-qe/integration1-iegh-tenant/integ-app-olhh/hacbs-test-project-c1k1:build-abb81-1699943742\\n\\nDetermine if Image Already Exists\\n\\ninit container prepare: \\n2023/11/14 06:36:50 Entrypoint initialization\\n\\ninit container place-scripts: \\n2023/11/14 06:36:51 Decoded script /tekton/scripts/script-0-m6zhw\\n\\ncontainer step-appstudio-summary: \\n\\nBuild Summary:\\n\\nBuild repository: https://github.com/redhat-appstudio-qe/hacbs-test-project?rev=34da5a8f51fba6a8b7ec75a727d3c72ebb5e1274\\n\\nEnd Summary\\n\",\n    }",
                            "failure_date": "2023-11-14T06:17:22Z"
                        }
                    ]
                },
                {
                    "name": "[It] [integration-service-suite Integration Service E2E tests] with happy path for general flow of Integration service when An snapshot of push event is created checks if an SnapshotEnvironmentBinding is created successfully [integration-service, HACBS]",
                    "test_case_impact": 2.7027027027027026,
                    "count": 1,
                    "messages": [
                        {
                            "job_id": "1724144866838974464",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2738/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724144866838974464",
                            "error_message": "Expected success, but got an error:\n    <context.deadlineExceededError>: \n    context deadline exceeded\n    {}",
                            "failure_date": "2023-11-13T19:19:03Z"
                        }
                    ]
                },
                {
                    "name": "[It] [integration-service-suite Integration Service E2E tests] with an integration test fail when An snapshot of push event is created checks if the global candidate is not updated [integration-service, HACBS]",
                    "test_case_impact": 2.7027027027027026,
                    "count": 1,
                    "messages": [
                        {
                            "job_id": "1724067359221616640",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2728/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724067359221616640",
                            "error_message": "Expected success, but got an error:\n    <context.deadlineExceededError>: \n    context deadline exceeded\n    {}",
                            "failure_date": "2023-11-13T14:11:04Z"
                        }
                    ]
                }
            ],
            "average_impact": 8.108108108108109
        },
        {
            "suite_name": "integration-service-suite Namespace-backed Environment (NBE) E2E tests",
            "status": "failed",
            "test_cases": [
                {
                    "name": "[It] [integration-service-suite Namespace-backed Environment (NBE) E2E tests] with happy path for Namespace-backed environments checks for deploymentTargetClaim after Ephemeral env has been created [integration-service, HACBS, namespace-backed-envs]",
                    "test_case_impact": 2.7027027027027026,
                    "count": 1,
                    "messages": [
                        {
                            "job_id": "1724272905153417216",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2741/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724272905153417216",
                            "error_message": "Timed out after 60.001s.\ntimed out checking DeploymentTargetClaim after Ephemeral Environment user-picked-environment-example-pass-with-env-sxgh-s7mdm was created \nExpected success, but got an error:\n    <*errors.errorString | 0xc00057d730>: \n    DeploymentTargetClaimPhase is not yet equal to the expected phase: Bound\n    {\n        s: \"DeploymentTargetClaimPhase is not yet equal to the expected phase: Bound\",\n    }",
                            "failure_date": "2023-11-14T03:47:50Z"
                        }
                    ]
                }
            ],
            "average_impact": 2.7027027027027026
        },
        {
            "suite_name": "integration-service-suite Status Reporting of Integration tests",
            "status": "failed",
            "test_cases": [
                {
                    "name": "[It] [integration-service-suite Status Reporting of Integration tests] with status reporting of Integration tests in CheckRuns when a new Component with specified custom branch is created eventually leads to the build PipelineRun's status reported at Checks tab [integration-service, HACBS, status-reporting, custom-branch]",
                    "test_case_impact": 2.7027027027027026,
                    "count": 1,
                    "messages": [
                        {
                            "job_id": "1724359083231809536",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2743/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724359083231809536",
                            "error_message": "the PR 78 in hacbs-test-project-integration repo doesn't contain the expected conclusion (success) of the CheckRun\nExpected\n    <string>: failure\nto equal\n    <string>: success",
                            "failure_date": "2023-11-14T09:30:17Z"
                        }
                    ]
                },
                {
                    "name": "[It] [integration-service-suite Status Reporting of Integration tests] with status reporting of Integration tests in CheckRuns when Integration PipelineRuns completes successfully eventually leads to the status reported at Checks tab for the successful Integration PipelineRun [integration-service, HACBS, status-reporting]",
                    "test_case_impact": 2.7027027027027026,
                    "count": 1,
                    "messages": [
                        {
                            "job_id": "1724410953295990784",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2728/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724410953295990784",
                            "error_message": "the PR 94 in hacbs-test-project-integration repo doesn't contain the expected conclusion (success) of the CheckRun\nExpected\n    <string>: \nto equal\n    <string>: success",
                            "failure_date": "2023-11-14T12:56:23Z"
                        }
                    ]
                }
            ],
            "average_impact": 5.405405405405405
        },
        {
            "suite_name": "rhtap-demo-suite",
            "status": "failed",
            "test_cases": [
                {
                    "name": "[It] [rhtap-demo-suite] DEVHAS-234: create an application with branch and context dir creates an environment [rhtap-demo]",
                    "test_case_impact": 2.7027027027027026,
                    "count": 1,
                    "messages": [
                        {
                            "job_id": "1724046292146982912",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2728/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724046292146982912",
                            "error_message": "Unexpected error:\n    <*url.Error | 0xc001f661e0>: \n    Post \"https://api-toolchain-host-operator.apps.rhtap-ocp-4-12-amd64-us-east-2-m8kt5.ci.stonesoupengineering.com/apis/appstudio.redhat.com/v1alpha1/namespaces/rhtap-demo-gihe-tenant/environments\": EOF\n    {\n        Op: \"Post\",\n        URL: \"https://api-toolchain-host-operator.apps.rhtap-ocp-4-12-amd64-us-east-2-m8kt5.ci.stonesoupengineering.com/apis/appstudio.redhat.com/v1alpha1/namespaces/rhtap-demo-gihe-tenant/environments\",\n        Err: <*errors.errorString | 0xc0001920f0>{s: \"EOF\"},\n    }\noccurred",
                            "failure_date": "2023-11-13T12:47:21Z"
                        }
                    ]
                },
                {
                    "name": "[It] [rhtap-demo-suite] Maven project - Simple and Advanced build RHTAP Advanced build test for rhtap-demo-component when SLSA level 3 customizable PipelineRun is created should eventually complete successfully [rhtap-demo]",
                    "test_case_impact": 2.7027027027027026,
                    "count": 1,
                    "messages": [
                        {
                            "job_id": "1724359083231809536",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2743/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724359083231809536",
                            "error_message": "Expected success, but got an error:\n    <*errors.errorString | 0xc00106a570>: \n    \n    {s: \"\"}",
                            "failure_date": "2023-11-14T09:30:17Z"
                        }
                    ]
                },
                {
                    "name": "[It] [rhtap-demo-suite] Application with a golang component with dockerfile but not devfile (private) waits for mc-golang-nodevfile component (private: true) pipeline to be finished [rhtap-demo]",
                    "test_case_impact": 2.7027027027027026,
                    "count": 1,
                    "messages": [
                        {
                            "job_id": "1724386610788700160",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2719/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724386610788700160",
                            "error_message": "Expected success, but got an error:\n    <*errors.errorString | 0xc001312fc0>: \n    \n    {s: \"\"}",
                            "failure_date": "2023-11-14T11:19:40Z"
                        }
                    ]
                },
                {
                    "name": "[It] [rhtap-demo-suite] multi-component scenario with all supported import components deploys component mc-three-scenarios successfully using gitops [rhtap-demo]",
                    "test_case_impact": 2.7027027027027026,
                    "count": 1,
                    "messages": [
                        {
                            "job_id": "1724386610788700160",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2719/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724386610788700160",
                            "error_message": "Timed out after 1500.000s.\ntimed out waiting for deployment of a component rhtap-demo-yigu-tenant/devfile-go-rhtap-three-component-scenarios-i8oj to become ready\nExpected success, but got an error:\n    <*errors.StatusError | 0xc000a09360>: \n    deployments.apps \"devfile-go-rhtap-three-component-scenarios-i8oj\" not found\n    {\n        ErrStatus: {\n            TypeMeta: {Kind: \"\", APIVersion: \"\"},\n            ListMeta: {\n                SelfLink: \"\",\n                ResourceVersion: \"\",\n                Continue: \"\",\n                RemainingItemCount: nil,\n            },\n            Status: \"Failure\",\n            Message: \"deployments.apps \\\"devfile-go-rhtap-three-component-scenarios-i8oj\\\" not found\",\n            Reason: \"NotFound\",\n            Details: {\n                Name: \"devfile-go-rhtap-three-component-scenarios-i8oj\",\n                Group: \"apps\",\n                Kind: \"deployments\",\n                UID: \"\",\n                Causes: nil,\n                RetryAfterSeconds: 0,\n            },\n            Code: 404,\n        },\n    }",
                            "failure_date": "2023-11-14T11:19:40Z"
                        }
                    ]
                }
            ],
            "average_impact": 10.81081081081081
        },
        {
            "suite_name": "spi-suite",
            "status": "failed",
            "test_cases": [
                {
                    "name": "[It] [spi-suite] SVPI-495 - Test automation to ensure that a user can't access and use secrets from another workspace checks that user A can access the SPIAccessToken A in workspace A [spi-suite, access-control]",
                    "test_case_impact": 2.7027027027027026,
                    "count": 1,
                    "messages": [
                        {
                            "job_id": "1724410953295990784",
                            "job_url": "https://prow.ci.openshift.org/view/gs/origin-ci-test/pr-logs/pull/redhat-appstudio_infra-deployments/2728/pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests/1724410953295990784",
                            "error_message": "Unexpected error:\n    <*url.Error | 0xc001752f30>: \n    Post \"https://api-toolchain-host-operator.apps.rhtap-ocp-4-12-amd64-us-east-2-chx2f.ci.stonesoupengineering.com/apis/appstudio.redhat.com/v1beta1/namespaces/spi-user-b-hxvs-tenant/spiaccesstokenbindings\": EOF\n    {\n        Op: \"Post\",\n        URL: \"https://api-toolchain-host-operator.apps.rhtap-ocp-4-12-amd64-us-east-2-chx2f.ci.stonesoupengineering.com/apis/appstudio.redhat.com/v1beta1/namespaces/spi-user-b-hxvs-tenant/spiaccesstokenbindings\",\n        Err: <*errors.errorString | 0xc0001980f0>{s: \"EOF\"},\n    }\noccurred",
                            "failure_date": "2023-11-14T12:56:23Z"
                        }
                    ]
                }
            ],
            "average_impact": 2.7027027027027026
        }
    ]
  }

  const mockImpact = [
      {
          "Date": "2023-11-10 23:59:59",
          "global_impact": 0
      },
      {
          "Date": "2023-11-11 23:59:59",
          "global_impact": 0
      },
      {
          "Date": "2023-11-12 23:59:59",
          "global_impact": 0
      },
      {
          "Date": "2023-11-13 23:59:59",
          "global_impact": 46.15384615384615
      },
      {
          "Date": "2023-11-14 23:59:59",
          "global_impact": 39.53488372093023
      },
      {
          "Date": "2023-11-15 23:59:59",
          "global_impact": 35.294117647058826
      },

     { "Date": "2023-11-16 23:59:59",
      "global_impact": 0
  },
  {
      "Date": "2023-11-17 23:59:59",
      "global_impact": 0
  },
  {
      "Date": "2023-11-18 23:59:59",
      "global_impact": 0
  },
  {
      "Date": "2023-11-19 23:59:59",
      "global_impact": 46.15384615384615
  },
  {
      "Date": "2023-11-20 23:59:59",
      "global_impact": 39.53488372093023
  },
  {
      "Date": "2023-11-21 23:59:59",
      "global_impact": 35.294117647058826
  }
  ]

  const countSuiteFailures = (suites) => {
    return suites.map((suite) => {
      const c = suite.test_cases.reduce(function (acc, obj) { return acc + obj.count; }, 0);
      return {suite_name: suite.suite_name, count: c}
    })
  }

  React.useEffect(() => {
    if(mockData && mockData.suites){
      setData(mockData.suites)
    }
    if(mockImpact){
      setBarData(mockImpact.map(impact => { impact.Date = impact.Date.split(' ')[0]; return impact;}))
    }
  }, []);

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
                  <DropdownBasic selected={selectedSuite} toggles={toggles} onSelect={onSuiteSelect}></DropdownBasic>
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
