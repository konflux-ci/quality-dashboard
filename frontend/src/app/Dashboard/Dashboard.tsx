import React, { useContext, useEffect, useState } from 'react';
import {
  Card,
  CardTitle,
  CardBody,
  Gallery,
  PageSection,
  PageSectionVariants,
  TextContent,
  Text,
  Title,
  Grid,
  GridItem,
  TitleSizes,
  DescriptionList, DescriptionListGroup, DescriptionListTerm, DescriptionListDescription,
  Drawer,
  DrawerPanelContent,
  DrawerContent,
  DrawerContentBody,
  DrawerHead,
  DrawerPanelBody,
  DrawerActions,
  DrawerCloseButton,
  List,
  ListItem,
} from '@patternfly/react-core';
import {
  TableComposable,
  Thead,
  Tr,
  Th,
  Tbody,
  Td,
  ActionsColumn,
  IAction,
  Caption,
  ThProps
} from '@patternfly/react-table';
import { ExternalLinkAltIcon } from '@patternfly/react-icons';
import { getVersion, getJiras } from '@app/utils/APIService';
import { Context } from "src/app/store/store";
import { RepositoriesTable } from '@app/Repositories/RepositoriesTable';
import { TableComponent } from '@app/Repositories/TableComponent';
import { parse } from 'postcss';


export const Dashboard = () => {
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
  const JIRA_CRITICAL = "Critical"
  const JIRA_BLOCKER = "Blocker"
  const JIRA_MAJOR = "Major"



  function computeJiraIssueCount(type) {
    try {
      return jiras.filter(j => j["fields"]["priority"]["name"] == type).length
    } catch (nullError) {
      return 0
    }
  }

  const issuesColumnNames = {
    issue_name: 'Issue',
    assignee: 'Assignee',
    dt_updated: 'Updated',
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

    if (dateTime["hour"] < 12) {
      dateTime['noon'] = 'AM'
    } else { dateTime['noon'] = 'PM' }

    return dateTime
  }


  const dt = null;
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
      {computeJiraIssueCount(jiraType) > 0 &&

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
        ))}
    </TableComposable>
  )
  
  // Sort helpers
  const [activeSortIndex, setActiveSortIndex] = React.useState<number | null>(null);
  const [activeSortDirection, setActiveSortDirection] = React.useState<'asc' | 'desc' | null>(null);
  const getSortableRowValues = (jiras): (string | number)[] => {
    const { issue_name, dt_updated, assignee, time_open } = jiras;
    return [issue_name, dt_updated, assignee, time_open];
  };

  let sortedJiras = jiras

  if (activeSortIndex !== null) {
    sortedJiras = jiras.sort((a, b) => {
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
    });
  }

  const getSortParams = (columnIndex: number): ThProps['sort'] => ({
    sortBy: {
      index: activeSortIndex as number,
      direction: activeSortDirection as any
    },
    onSort: (_event, index, direction) => {
      setActiveSortIndex(index);
      setActiveSortDirection(direction);
    },
    columnIndex
  });

  const panelContent = (
    <DrawerPanelContent isResizable defaultSize={'15vw'} minSize={'35vw'}>
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

  const [dashboardVersion, setVersion] = useState('unknown')
  const { state, dispatch } = useContext(Context) // required to access the global state
  useEffect(() => {
    getVersion().then((res) => { // making the api call here
      if (res.code === 200) {
        const result = res.data;
        dispatch({ type: "SET_Version", data: result['version'] });
        // not really required to store it in the global state , just added it to make it better understandable
        setVersion(result['version'])
      } else {
        dispatch({ type: "SET_ERROR", data: res });
      }
    });
  }, [dashboardVersion, setVersion, dispatch])

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
  }, [])

  return (
    <React.Fragment>
      <PageSection style={{
        minHeight: "12%",
        background: "url(https://console.redhat.com/apps/frontend-assets/background-images/new-landing-page/estate_section_banner.svg)",
        backgroundSize: "cover",
        backgroundColor: "black",
        opacity: '0.9'
      }} variant={PageSectionVariants.light}>
        <TextContent style={{ color: "white" }}>
          <Text component="h2">Red Hat App Studio Quality Dashboard</Text>
          <Text component="p">This is a demo that show app studio quality status.</Text>
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
                      <DescriptionListTerm>Quality Dashboard version</DescriptionListTerm>
                      <DescriptionListDescription>
                        <span>{dashboardVersion}</span>
                      </DescriptionListDescription>
                    </DescriptionListGroup>
                    <DescriptionListGroup>
                      <DescriptionListTerm>Staging Version</DescriptionListTerm>
                      <DescriptionListDescription>
                        <span>Unknown Version</span>
                      </DescriptionListDescription>
                    </DescriptionListGroup>
                    <DescriptionListGroup>
                      <DescriptionListTerm>Production Version</DescriptionListTerm>
                      <DescriptionListDescription>Unknown Version</DescriptionListDescription>
                    </DescriptionListGroup>
                    <DescriptionListGroup>
                      <DescriptionListTerm>Github Organization</DescriptionListTerm>
                      <a href="https://github.com/redhat-appstudio">redhat-appstudio <ExternalLinkAltIcon ></ExternalLinkAltIcon></a>
                    </DescriptionListGroup>
                  </DescriptionList>
                </CardBody>
              </Card>
              <Card isRounded isCompact style={{ width: "65%" }}>
                <CardTitle>
                  <Title headingLevel="h2" size="xl" >
                    Active jira Issues
                  </Title>
                </CardTitle>
                <Grid md={4} style={{ margin: "auto 5px" }}>
                  <GridItem style={{ margin: "5px", cursor: "pointer" }} aria-expanded={isExpanded} onClick={event => showJiras(JIRA_BLOCKER)}>
                    <Card>
                      <CardBody style={{ display: "flex", alignItems: "center", justifyContent: "center", minHeight: "200px", margin: "auto 5px", textAlign: "center" }}>
                        <div>
                          <Title headingLevel="h1" size={TitleSizes['4xl']}>{computeJiraIssueCount(JIRA_BLOCKER)}</Title>
                          <p>Blocker</p>
                        </div>
                      </CardBody>
                    </Card>
                  </GridItem>

                  <GridItem style={{ margin: "5px", cursor: "pointer" }} aria-expanded={isExpanded} onClick={event => showJiras(JIRA_CRITICAL)}>
                    <Card>
                      <CardBody style={{ display: "flex", alignItems: "center", justifyContent: "center", minHeight: "200px", margin: "auto 5px", textAlign: "center" }}>
                        <div>
                          <Title headingLevel="h1" size={TitleSizes['4xl']}>{computeJiraIssueCount(JIRA_CRITICAL)}</Title>
                          <p>Critical</p>
                        </div>
                      </CardBody>
                    </Card>
                  </GridItem>

                  <GridItem style={{ margin: "5px", cursor: "pointer" }} aria-expanded={isExpanded} onClick={event => showJiras(JIRA_MAJOR)}>
                    <Card>
                      <CardBody style={{ display: "flex", alignItems: "center", justifyContent: "center", minHeight: "200px", margin: "auto 5px", textAlign: "center" }}>
                        <div>
                          <Title headingLevel="h1" size={TitleSizes['4xl']}>{computeJiraIssueCount(JIRA_MAJOR)}</Title>
                          <p>Major</p>
                        </div>
                      </CardBody>
                    </Card>
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
