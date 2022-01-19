import React from 'react';
import {  PageSectionVariants, TextContent, Text, PageSection} from '@patternfly/react-core';
import { FormModal } from './CreateRepository';
import { TableComponent } from './TableComponent';

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
        <FormModal/>
      </PageSection>
      <PageSection>
        <TableComponent/>
      </PageSection>
    </React.Fragment>
  );
};
