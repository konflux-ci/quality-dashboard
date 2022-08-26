import React, { useContext, useEffect, useState } from 'react';
import {
  Card,
  CardTitle,
  CardBody,
  Gallery,
  PageSection,
  Title,
  DescriptionList, DescriptionListGroup, DescriptionListTerm, 
  DescriptionListDescription, TextContent, Text, PageSectionVariants
} from '@patternfly/react-core';
import { CheckCircleIcon, ExclamationCircleIcon } from '@patternfly/react-icons';
import { getVersion } from '@app/utils/APIService';
import { Context } from "src/app/store/store";
import { RepositoriesTable } from '@app/Repositories/RepositoriesTable';

export const Dashboard = () => {
  const [dashboardVersion, setVersion] = useState('unknown')
  const [serverAvailable, setServerAvailable] = useState<boolean>(false)
  const {state, dispatch} = useContext(Context) // required to access the global state
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
            <Text component="p">This is a demo that show app studio quality status.</Text>
          </TextContent>
        </PageSection>
        <PageSection>
          <Gallery hasGutter style={{ display:"flex" }}>
            <Card style={{width: "35%"}}>
              <CardTitle>
                <Title headingLevel="h1" size="xl">
                  Red Hat Quality Studio Status
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
            <Card style={{width: "65%"}}>
              <CardTitle>
                <Title headingLevel="h2" size="xl">
                  Tests Summary
                </Title>
              </CardTitle>
            </Card>
          </Gallery>
          </PageSection>
        <PageSection style={{
        minHeight : "12%"
      }}>
        <RepositoriesTable showTableToolbar={true} showCoverage={true} showDiscription={true} enableFiltersOnTheseColumns={['git_organization']}></RepositoriesTable>
      </PageSection>
    </React.Fragment>
  );
}
