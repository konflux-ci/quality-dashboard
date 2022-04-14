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
  DescriptionList, DescriptionListGroup, DescriptionListTerm, DescriptionListDescription
} from '@patternfly/react-core';
import { TableComponent } from '@app/Repositories/TableComponent';
import { ExternalLinkAltIcon } from '@patternfly/react-icons';
import { getVersion } from '@app/utils/APIService';
import { Context } from "src/app/store/store";

export const Dashboard = () => {
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
          <PageSection >
          <Gallery hasGutter style={{ display:"flex" }}>
            <Card style={{width: "35%"}}>
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
        <TableComponent showCoverage showDiscription={false}></TableComponent>
      </PageSection>
    </React.Fragment>
  );
}
