import React, { useState, useContext } from 'react';
import {
  PageSection,
  PageSectionVariants,
  AlertVariant,
  Button,
  Toolbar,
  ToolbarContent,
  ToolbarItem,
  ToolbarGroup,
  Title
} from '@patternfly/react-core';
import { PlusIcon } from '@patternfly/react-icons/dist/esm/icons';
import { ReactReduxContext, useSelector } from 'react-redux';

interface AlertInfo {
  title: string;
  variant: AlertVariant;
  key: string;
}

import { Table, TableHeader, TableBody, TableProps } from '@patternfly/react-table';
import { TeamForm } from './TeamForm';



export const TeamsWizard = () => {
  const [isOpen, setOpen] = useState<boolean>(false);

  const handleModalToggle = () => {
    setOpen(!isOpen)
  };

  function handleChange(event, open) {
    setOpen(open)
  }

  return (
    <React.Fragment>
      <PageSection style={{ backgroundColor: 'white' }} variant={PageSectionVariants.light}>
        <Toolbar id="toolbar-items">
          <ToolbarContent>
            <ToolbarGroup variant="filter-group" alignment={{ default: 'alignLeft' }}>
              <Title headingLevel="h2" size="3xl">Teams</Title>
            </ToolbarGroup>
            <ToolbarGroup variant="filter-group" alignment={{ default: 'alignRight' }}>
              <ToolbarItem>
                <Button onClick={handleModalToggle} type="button" variant="primary"> <PlusIcon></PlusIcon> Add Team </Button>
              </ToolbarItem>
            </ToolbarGroup>
          </ToolbarContent>
        </Toolbar>
        <TeamsTable></TeamsTable>
        <TeamForm
          isOpen={isOpen}
          handleChange={(event, open) => handleChange(event, open)}
        ></TeamForm>
      </PageSection>
    </React.Fragment>
  );
};

export const TeamsTable: React.FunctionComponent = () => {
  // In real usage, this data would come from some external source like an API via props.

  const { store } = useContext(ReactReduxContext);
  const state = store.getState();
  const columns: TableProps['cells'] = ['Name', 'Description'];

  let currentTeamsAvailable = useSelector((state: any) => state.teams.TeamsAvailable);
  //store.subscribe(currentTeamsAvailable = useSelector((state:any) => state.teams.TeamsAvailable));

  const rows: TableProps['rows'] = currentTeamsAvailable.map(team => [
    team.team_name,
    team.description
  ]);

  return (
    <React.Fragment>
      <Table
        aria-label="Teams Table"
        cells={columns}
        rows={rows}
      >
        <TableHeader />
        <TableBody />
      </Table>
    </React.Fragment>
  );
};