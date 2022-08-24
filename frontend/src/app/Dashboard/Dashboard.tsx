import React, { useContext, useEffect, useState } from 'react';
import {
  Card,
  CardTitle,
  CardBody,
  Gallery,
  PageSection,
  Title,
  DescriptionList, DescriptionListGroup, DescriptionListTerm, DescriptionListDescription, Button, ButtonVariant
} from '@patternfly/react-core';
import { ExternalLinkAltIcon, PlusIcon } from '@patternfly/react-icons';
import { getVersion } from '@app/utils/APIService';
import { Context } from "src/app/store/store";
import { RepositoriesTable } from '@app/Repositories/RepositoriesTable';


export const Dashboard = () => {
  const [dashboardVersion, setVersion] = useState('unknown')
  const {state, dispatch} = useContext(Context) // required to access the global state
  useEffect(()=> {
    getVersion().then((res) => { // making the api call here
      if(res.code === 200){
          const result = res.data;
          dispatch({ type: "SET_Version", data: result['serverAPIVersion'] });
          // not really required to store it in the global state , just added it to make it better understandable
          setVersion(result['serverAPIVersion'])
      } else {
          dispatch({ type: "SET_ERROR", data: res });
      }
    });
  }, [dashboardVersion, setVersion, dispatch])

  return (
    <React.Fragment>
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
                      <span>Unknown Version</span>
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                  <DescriptionListGroup>
                    <DescriptionListTerm>Database Status</DescriptionListTerm>
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
        <RepositoriesTable showTableToolbar={true} showCoverage={true} showDiscription={true} enableFiltersOnTheseColumns={['git_organization']}></RepositoriesTable>
      </PageSection>
    </React.Fragment>
  );
}
