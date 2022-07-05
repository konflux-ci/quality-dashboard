import React, { useState } from 'react';
import {  PageSectionVariants, TextContent, Text, PageSection} from '@patternfly/react-core';
import { RepositoriesTable } from './RepositoriesTable';

export interface Coverage {
  coverage_percentage: number
}

export interface Repository {
  git_organization: string;
  repository_name: string;
  git_url: string;
  description: string;
  code_coverage: Coverage;
}

export const Repositories: React.FunctionComponent = () => {
  return (
    <React.Fragment>
      <PageSection style={{
          backgroundSize: "cover",
          backgroundColor : "black",
        }} variant={PageSectionVariants.darker}>
        <TextContent style={{color: "white"}}>
          <Text component="h2">Red Hat App Studio Quality Dashboard</Text>
          <Text component="p">Repositories list</Text>
        </TextContent>
      </PageSection>
      <PageSection>
        <RepositoriesTable showTableToolbar={true} showCoverage={false} showDiscription={true} enableFiltersOnTheseColumns={['repository_name', 'git_organization']}></RepositoriesTable>
      </PageSection>
    </React.Fragment>
  );
};
