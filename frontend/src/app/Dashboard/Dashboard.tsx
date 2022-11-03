import React, { useContext, useEffect, useState } from 'react';
import {
  Card,
  CardTitle,
  CardBody,
  Gallery,
  PageSection,
  Title,
  DescriptionList, DescriptionListGroup, DescriptionListTerm, 
  DescriptionListDescription, TextContent, Text, PageSectionVariants,
  Grid,
  GridItem,
  TitleSizes,
  Drawer,
  DrawerPanelContent,
  DrawerContent,
  DrawerContentBody,
  DrawerHead,
  DrawerPanelBody,
  DrawerActions,
  DrawerCloseButton,
} from '@patternfly/react-core';
import {
  TableComposable,
  Thead,
  Tr,
  Th,
  Tbody,
  Td,
  ThProps
} from '@patternfly/react-table';
import { ChartDonut, ChartPie, ChartThemeColor } from '@patternfly/react-charts';
import { ExternalLinkAltIcon, } from '@patternfly/react-icons';
import { getVersion, getJiras } from '@app/utils/APIService';
import { Context } from "src/app/store/store";
import { RepositoriesTable } from '@app/Repositories/RepositoriesTable';
import { TableComponent } from '@app/Repositories/TableComponent';
import { parse } from 'postcss';
import { string } from 'prop-types';
import { toInteger } from 'lodash';
import { CheckCircleIcon, ExclamationCircleIcon } from '@patternfly/react-icons';
import { ReactReduxContext } from 'react-redux';

export const Dashboard = () => {

  /*
    ALL JIRA RELATED STUFF
  */
  const [isExpanded, setIsExpanded] = React.useState(false);
  const drawerRef = React.useRef<HTMLDivElement>();

  const onExpand = () => {
    drawerRef.current && drawerRef.current.focus();
  };

  const showJiras = (issueType) => {
    setIsExpanded(!isExpanded);
    setJiraType(issueType)
  };

  const onCloseClick = () => {
    setIsExpanded(false);
  };

  const [jiraType, setJiraType] = useState("");
  const [jiras, setJiras] = useState([]);
  const JIRA_ALL = "All"
  const JIRA_CRITICAL = "Critical"
  const JIRA_BLOCKER = "Blocker"
  const JIRA_MAJOR = "Major"

  useEffect(() => {
    getJiras().then((res) => {
      if (res.code === 200) {
        const result = res.data;
        dispatch({ type: "SET_JIRAS", data: result });
        setJiras(result)
      } else {
        dispatch({ type: "SET_ERROR", data: res });
      }
    })
  }, [jiras, setJiras]);

  function computeJiraIssueCount(type) {
    try {
      if (type == JIRA_ALL) {
        return computeJiraIssueCount(JIRA_BLOCKER) + computeJiraIssueCount(JIRA_CRITICAL) + computeJiraIssueCount(JIRA_MAJOR)
      }
      return jiras.filter(j => j["fields"]["priority"]["name"] == type).length
    } catch (nullError) {
      return 0
    }
  }

  const issuesColumnNames = {
    issue_name: 'Issue',
    assignee: 'Assignee',
    dt_updated: 'Date',
    time_open: 'Age'
  };

  const monthShort = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
  function parseDate(dt) {
    const dateTime = []
    dateTime['year'] = parseInt(dt.slice(0, 4))
    dateTime['month'] = monthShort[parseInt(dt.slice(5, 7)) - 1]
    dateTime['day'] = parseInt(dt.slice(8, 10))
    dateTime['hour'] = parseInt(dt.slice(11, 13))
    dateTime['minute'] = parseInt(dt.slice(14, 16))

    if (dateTime['hour'] < 12) {
      dateTime['noon'] = 'AM'
    } else { dateTime['noon'] = 'PM' }

    return dateTime
  }

  const totalIssues = computeJiraIssueCount(JIRA_BLOCKER) + computeJiraIssueCount(JIRA_CRITICAL) + computeJiraIssueCount(JIRA_MAJOR)
  const allIssueChart = () => (
    <div style={{ height: '150px', width: '150px' }}>
      <ChartDonut
        ariaDesc="Jira Issues All"
        constrainToVisibleArea
        data={[{ x: 'Blocker', y: computeJiraIssueCount(JIRA_BLOCKER) }, { x: 'Critical', y: computeJiraIssueCount(JIRA_CRITICAL) }, { x: 'Major', y: computeJiraIssueCount(JIRA_MAJOR) }]}
        height={150}
        labels={({ datum }) => `${datum.x}: ${datum.y}/${totalIssues}`}
        title={totalIssues.toString()}
        themeColor={ChartThemeColor.blue}
        width={150}
      />
    </div>

  )

  const issueChart = (jiraType) => (
    <div style={{ height: '150px', width: '150px' }}>
      <ChartDonut
        ariaDesc="Issues by type"
        constrainToVisibleArea
        data={[{ x: 'Issues', y: computeJiraIssueCount(jiraType) }, { x: 'Whitespace', y: 15 - computeJiraIssueCount(jiraType) }]}
        height={150}
        labels={({ datum }) => null}
        title={computeJiraIssueCount(jiraType).toString()}
        themeColor={ChartThemeColor.blue}
        width={150}
      />
    </div>

  )

  const dt = null;
  let visibleJiras = null;
  const JiraIssuesList = () => (
    <TableComposable aria-label="Jiras table">
      <Thead>
        <Tr>
          <Th>{issuesColumnNames.issue_name}</Th>
          <Th>{issuesColumnNames.dt_updated}</Th>
          <Th>{issuesColumnNames.assignee}</Th>
        </Tr>
      </Thead>
      {computeJiraIssueCount(jiraType) == 0 &&
        <div style={{ textAlign: "center", margin: "10px auto", minHeight: "500px" }}><i>No issues here</i></div>
      }
      {computeJiraIssueCount(jiraType) > 0 && jiraType == JIRA_ALL &&
        jiras.map(j => (
          <Tbody style={{ marginTop: "5px" }}>
            <Tr>
              <Td>
                <div>
                  <strong style={{ textDecoration: "underline", color: "blue" }}><a href={`https://issues.redhat.com/browse/${j["key"]}`}>{j["key"]}</a></strong>
                  : &nbsp;
                </div>
                <div>{j["fields"]["summary"]}</div>
              </Td>
              <Td><div>{parseDate(j['fields']['updated'])['month']} {parseDate(j['fields']['updated'])['day']}, {parseDate(j['fields']['updated'])['year']}</div>
                <div>{parseDate(j['fields']['updated'])['hour']}:{parseDate(j['fields']['updated'])['minute']}</div>
              </Td>
              <Td>{j["fields"]["assignee"]["displayName"]}</Td>
            </Tr>
          </Tbody>
        ))
      }
      {computeJiraIssueCount(jiraType) > 0 && jiraType != JIRA_ALL &&
        jiras.filter(j => j["fields"]["priority"]["name"] == jiraType).map(j => (
          <Tbody style={{ marginTop: "5px" }}>
            <Tr>
              <Td>
                <div>
                  <strong style={{ textDecoration: "underline", color: "blue" }}><a href={`https://issues.redhat.com/browse/${j["key"]}`}>{j["key"]}</a></strong>
                  : &nbsp;
                </div>
                <div>{j["fields"]["summary"]}</div>
              </Td>
              <Td><div>{parseDate(j['fields']['updated'])['month']} {parseDate(j['fields']['updated'])['day']}, {parseDate(j['fields']['updated'])['year']}</div>
                <div>{parseDate(j['fields']['updated'])['hour']}:{parseDate(j['fields']['updated'])['minute']}</div>
              </Td>
              <Td>{j["fields"]["assignee"]["displayName"]}</Td>
            </Tr>
          </Tbody>
        ))
      }

    </TableComposable>
  )

  const panelContent = (
    <DrawerPanelContent isResizable defaultSize={'15vw'} minSize={'45vw'}>
      <DrawerHead>
        <DrawerActions>
          <DrawerCloseButton onClick={onCloseClick} />
        </DrawerActions>
      </DrawerHead>
      <DrawerPanelBody>
        <div>
          <Title headingLevel="h1" size="xl" style={{ textTransform: "uppercase", marginBottom: "10px" }}>{jiraType} Issues</Title>
          <JiraIssuesList></JiraIssuesList>
        </div>
      </DrawerPanelBody>
    </DrawerPanelContent>
  );

  /*
    SERVER INFO AND REPOSITORY TABLE METHODS
  */
  const [dashboardVersion, setVersion] = useState('unknown')
  const [serverAvailable, setServerAvailable] = useState<boolean>(false)
  const { store } = useContext(ReactReduxContext);  
  const state = store.getState();
  const dispatch = store.dispatch;  

  useEffect(()=> {
    getVersion().then((res) => { // making the api call here
      if(res.code === 200){
          const result = res.data;
          dispatch({ type: "SET_Version", data: result['serverAPIVersion'] });
          // not really required to store it in the global state , just added it to make it better understandable
          setVersion(result['serverAPIVersion'])
          setServerAvailable(true)
      } else {
        setServerAvailable(false)
        dispatch({ type: "SET_ERROR", data: res });
      }
    });
  }, [dashboardVersion, setVersion, dispatch])

  return (
    <React.Fragment>
        <PageSection style={{
          minHeight : "12%",
          background:"url(https://console.redhat.com/apps/frontend-assets/background-images/new-landing-page/estate_section_banner.svg)",
          backgroundSize: "cover",
          backgroundColor : "black",
          opacity: '0.9'
        }} variant={PageSectionVariants.light}>
          <TextContent style={{color: "white"}}>
            <Text component="h2">Red Hat App Studio Quality Dashboard</Text>
            <Text component="p">Global status about teams components.</Text>
          </TextContent>
        </PageSection>
        <Drawer isExpanded={isExpanded}>
        <DrawerContent panelContent={panelContent} className={'pf-m-no-background'}>
          <PageSection>
            <Gallery hasGutter style={{ display: "flex" }}>
              <Card isRounded style={{ width: "35%" }}>
                <CardTitle>
                  <Title headingLevel="h1" size="xl">
                    Red Hat App Studio Details
                  </Title>
                </CardTitle>
                <CardBody>
                  <DescriptionList>
                  <DescriptionListGroup>
                      <DescriptionListTerm>Quality Studio version</DescriptionListTerm>
                      <DescriptionListDescription>
                        <span>{dashboardVersion}</span>
                      </DescriptionListDescription>
                    </DescriptionListGroup>
                    <DescriptionListGroup>
                      <DescriptionListTerm>Server API Status</DescriptionListTerm>
                      <DescriptionListDescription>
                        { serverAvailable &&  <span style={{color: "darkgreen", verticalAlign: "middle", lineHeight: "2em", fontWeight: 500}}> <CheckCircleIcon size={'sm'} ></CheckCircleIcon> OK </span> }
                        {!serverAvailable &&  <span style={{color: "darkred", verticalAlign: "middle", lineHeight: "2em", fontWeight: 500}}> <ExclamationCircleIcon size={'sm'} ></ExclamationCircleIcon> DOWN </span> }
                      </DescriptionListDescription>
                    </DescriptionListGroup>
                    <DescriptionListGroup>
                      <DescriptionListTerm>Database Status</DescriptionListTerm>
                      <DescriptionListDescription>
                        { serverAvailable &&  <span style={{color: "darkgreen", verticalAlign: "middle", lineHeight: "2em", fontWeight: 500}}> <CheckCircleIcon size={'sm'} ></CheckCircleIcon> OK </span> }
                        {!serverAvailable &&  <span style={{color: "darkred", verticalAlign: "middle", lineHeight: "2em", fontWeight: 500}}> <ExclamationCircleIcon size={'sm'} ></ExclamationCircleIcon> DOWN </span> }
                      </DescriptionListDescription>
                    </DescriptionListGroup>
                  </DescriptionList>
                </CardBody>
              </Card>
              <Card isRounded isCompact style={{ width: "65%", padding: "5px" }}>
                <CardTitle>
                  <Title headingLevel="h2" size="xl" > Jira Issues blocking e2e-tests </Title>
                </CardTitle>
                <Grid hasGutter span={3}>
                  <GridItem aria-expanded={isExpanded} onClick={event => showJiras(JIRA_ALL)}>
                    <CardBody style={{ display: "flex", alignItems: "center", justifyContent: "center", minHeight: "150px", margin: "auto 5px", textAlign: "center" }}>
                      <div style={{ cursor: "pointer" }} aria-expanded={isExpanded} onClick={event => showJiras(JIRA_ALL)}>
                        <Title headingLevel="h2" size={TitleSizes['lg']}>All</Title>
                        {allIssueChart()}
                      </div>
                    </CardBody>
                  </GridItem>
                  <GridItem style={{ cursor: "pointer" }} aria-expanded={isExpanded} onClick={event => showJiras(JIRA_BLOCKER)}>
                    <CardBody style={{ display: "flex", alignItems: "center", justifyContent: "center", minHeight: "200px", margin: "auto 5px", textAlign: "center" }}>
                      <div>
                        <Title headingLevel="h2" size={TitleSizes['lg']}>Blocker</Title>
                        {issueChart(JIRA_BLOCKER)}
                      </div>
                      </CardBody>
                  </GridItem>
                    <GridItem style={{ cursor: "pointer" }} aria-expanded={isExpanded} onClick={event => showJiras(JIRA_CRITICAL)}>
                      <CardBody style={{ display: "flex", alignItems: "center", justifyContent: "center", minHeight: "200px", margin: "auto 5px", textAlign: "center" }}>
                        <div>
                          <Title headingLevel="h2" size={TitleSizes['lg']}>Critical</Title>
                          {issueChart(JIRA_CRITICAL)}
                        </div>
                      </CardBody>
                      </GridItem>
                    <GridItem style={{ cursor: "pointer" }} aria-expanded={isExpanded} onClick={event => showJiras(JIRA_MAJOR)}>
                      <CardBody style={{ display: "flex", alignItems: "center", justifyContent: "center", minHeight: "200px", margin: "auto 5px", textAlign: "center" }}>
                        <div>
                          <Title headingLevel="h2" size={TitleSizes['lg']}>Major</Title>
                          {issueChart(JIRA_MAJOR)}
                        </div>
                      </CardBody>
                    </GridItem>
                </Grid>
              </Card>
            </Gallery>
          </PageSection>
          <PageSection style={{
            minHeight: "12%"
          }}>
            <RepositoriesTable showTableToolbar={true} showCoverage={true} showDiscription={false} enableFiltersOnTheseColumns={['git_organization']}></RepositoriesTable>
            <React.Fragment>
            </React.Fragment>
          </PageSection>
          <PageSection padding={{ default: 'noPadding' }}>
            <DrawerContentBody hasPadding></DrawerContentBody>
          </PageSection>
        </DrawerContent>
      </Drawer>
    </React.Fragment>
  );
}