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
import { ExternalLinkAltIcon } from '@patternfly/react-icons';
import { getVersion, getJiras } from '@app/utils/APIService';
import { Context } from "src/app/store/store";
import { RepositoriesTable } from '@app/Repositories/RepositoriesTable';
import { TableComponent } from '@app/Repositories/TableComponent';


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

  function computeJiraIssueCount(type){
    try { return jiras.filter(j => j["fields"]["priority"]["name"] == type).length
    } catch(nullError) { 
      return 0 
    }
  }

  const JiraIssuesList = () => (
    <List isPlain isBordered>
      { computeJiraIssueCount(jiraType) == 0 && 
        <div style={{textAlign: "center", margin: "10px auto", minHeight:"500px"}}><i>No issues here</i></div>
      }
      { computeJiraIssueCount(jiraType) > 0 && 
        jiras.filter(j => j["fields"]["priority"]["name"] == jiraType).map(j => (
          <ListItem style={{marginTop: "5px"}}>
            <strong style={{textDecoration: "underline", color: "blue"}}><a href={`https://issues.redhat.com/browse/${j["key"]}`}>{j["key"]}</a></strong>
            : &nbsp;
          {j["fields"]["summary"]}
        </ListItem>
      ))}
    </List>
  )

  const panelContent = (
    <DrawerPanelContent isResizable defaultSize={'15vw'} minSize={'35vw'}>
      <DrawerHead>
        <DrawerActions>
          <DrawerCloseButton onClick={onCloseClick} />
        </DrawerActions>
      </DrawerHead>
      <DrawerPanelBody>
        <div>
          <Title headingLevel="h1" size="xl" style={{textTransform: "uppercase", marginBottom: "10px"}}>{jiraType}</Title>
          <JiraIssuesList></JiraIssuesList>
        </div>
      </DrawerPanelBody>
    </DrawerPanelContent>
  );

  const [dashboardVersion, setVersion] = useState('unknown')
  const {state, dispatch} = useContext(Context) // required to access the global state
  useEffect(()=> {
    getVersion().then((res) => { // making the api call here
      if(res.code === 200){
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
    getJiras().then((res) =>{
      if(res.code === 200){
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
          minHeight : "12%",
          background:"url(https://console.redhat.com/apps/frontend-assets/background-images/new-landing-page/estate_section_banner.svg)",
          backgroundSize: "cover",
          backgroundColor : "black",
          opacity: '0.9'
        }} variant={PageSectionVariants.light}>
          <TextContent style={{color: "white"}}>
            <Text component="h2">Red Hat App Studio Quality Dashboard</Text>
            <Text component="p">This is a demo that show app studio quality status.</Text>
          </TextContent>
        </PageSection>
          <Drawer isExpanded={isExpanded}>
            <DrawerContent panelContent={panelContent} className={'pf-m-no-background'}>

          <PageSection>
          <Gallery hasGutter style={{ display:"flex" }}>
            <Card isRounded style={{width: "35%"}}>
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
            <Card isRounded isCompact style={{width: "65%"}}>
              <CardTitle>
                <Title headingLevel="h2" size="xl" >
                Active jira Issues
                </Title>
              </CardTitle>
              <Grid md={4} style={{margin: "auto 5px"}}>
                <GridItem style={{margin: "5px", cursor: "pointer"}} aria-expanded={isExpanded} onClick={event => showJiras(JIRA_BLOCKER)}>
                  <Card>
                    <CardBody style={{display: "flex", alignItems: "center", justifyContent: "center", minHeight: "200px", margin: "auto 5px", textAlign: "center"}}>
                      <div>
                        <Title headingLevel="h1" size={TitleSizes['4xl']}>{ computeJiraIssueCount(JIRA_BLOCKER) }</Title>
                        <p>Blocker</p>
                      </div>  
                    </CardBody>
                  </Card>
                </GridItem>

                <GridItem style={{margin: "5px", cursor: "pointer"}} aria-expanded={isExpanded} onClick={event => showJiras(JIRA_CRITICAL)}>
                  <Card>
                    <CardBody style={{display: "flex", alignItems: "center", justifyContent: "center", minHeight: "200px", margin: "auto 5px", textAlign: "center"}}>
                      <div>
                        <Title headingLevel="h1" size={TitleSizes['4xl']}>{ computeJiraIssueCount(JIRA_CRITICAL) }</Title>
                        <p>Critical</p>
                      </div>
                    </CardBody>
                  </Card>
                </GridItem>

                <GridItem style={{margin: "5px", cursor: "pointer"}} aria-expanded={isExpanded} onClick={event => showJiras(JIRA_MAJOR)}>
                  <Card>
                    <CardBody style={{display: "flex", alignItems: "center", justifyContent: "center", minHeight: "200px", margin: "auto 5px", textAlign: "center"}}>
                      <div>
                        <Title headingLevel="h1" size={TitleSizes['4xl']}>{ computeJiraIssueCount(JIRA_MAJOR) }</Title>
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
              minHeight : "12%"
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
