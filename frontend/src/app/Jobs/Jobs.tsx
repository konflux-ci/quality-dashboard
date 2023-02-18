/* eslint-disable react-hooks/exhaustive-deps */
/* eslint-disable @typescript-eslint/no-explicit-any */
import React, { useContext, useEffect, useState } from 'react';
import { Button, ContextSelector, ContextSelectorItem, EmptyState, EmptyStateIcon, EmptyStateVariant, PageSection, Title } from '@patternfly/react-core';
import { getRepositories, getWorkflowByRepositoryName } from '@app/utils/APIService';
import { Caption, TableComposable, Tbody, Td, Th, Thead, Tr } from '@patternfly/react-table';
import { ReactReduxContext, useSelector } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { CopyIcon, ExclamationCircleIcon, ExternalLinkAltIcon } from '@patternfly/react-icons';
import { isValidTeam } from '@app/utils/utils';

interface Repository {
  git_organization: string;
  repository_name: string;
  git_url: string
  description: string
}

interface Workflows {
  workflow_name: string;
  badge_url: string;
  job: string;
  html_url: string
}

const columnNames = {
  name: 'Name',
  html_url: 'Job URL',
  job: 'Job State',
  badge: 'Last Execution'
};


export const JobsComponent: React.FunctionComponent = () => {
  const [isOpen, setOpen] = useState(false);
  const [searchValue, setSearchValue] = useState('');

  const { store } = useContext(ReactReduxContext);
  const state = store.getState();
  const dispatch = store.dispatch;

  const [workflows, setWorkflows] = useState([]);

  const [repositories, setAllReps] = useState<any>(state.repos.Allrepositories);

  const [, setFilteredItems] = useState(repositories);
  const [repos, setRepositories] = useState([]);

  const currentTeam = useSelector((state: any) => state.teams.Team);

  let firstItemText = "default"
  const [selected, setSelected] = useState(firstItemText);
  const [org, setOrg] = useState('');

  const history = useHistory();
  const params = new URLSearchParams(window.location.search)

  useEffect(() => {
    clearAll()

    const repository = params.get("repository")
    const organization = params.get("organization")
    const team = params.get("team")

    getRepositories(5, state.teams.Team).then((res) => {
      if (res.code === 200) {
        const result = res.data;
        res.all.sort((a, b) => (a.repository_name < b.repository_name ? -1 : 1));
        setAllReps(res.all)
        dispatch({ type: "SET_REPOSITORIES", data: result });
        dispatch({ type: "SET_REPOSITORIES_ALL", data: res.all });

        if (res.all.length < 1 && (team == state.teams.Team || team == null)) {
          history.push('/ci/jobs?team=' + currentTeam)
        }

        if (res.all.length > 0 && (team == state.teams.Team || team == null)) {
          if (repository == null || organization == null) { // first click on GitHub Actions or team
            setSelected(res.all[0].repository_name)
            setOrg(res.all[0].git_organization)
            getWorkflows(res.all[0].repository_name)
            history.push('/ci/jobs?team=' + currentTeam + '&organization=' + res.all[0].git_organization + '&repository=' + res.all[0].repository_name)
          } else {
            setSelected(repository)
            setOrg(organization)
            getWorkflows(repository)
            history.push('/ci/jobs?team=' + currentTeam + '&organization=' + organization + '&repository=' + repository)
          }
        }

      } else {
        dispatch({ type: "SET_ERROR", data: res });
        clearAll()
      }
    });
  }, [repos, setRepositories, setAllReps, dispatch, state.teams.Team, currentTeam])


  function onToggle(_event: any, isOpen: boolean) {
    setOpen(isOpen);
  }

  function clearAll() {
    setWorkflows([])
    setSearchValue('')
    setSelected('default')
    setOrg('')
  }

  function getWorkflows(repo) {
    getWorkflowByRepositoryName(repo).then((res) => {
      if (res.code === 200) {
        const result = res.data.sort((a, b) => (a.workflow_name < b.workflow_name ? -1 : 1));
        setWorkflows(result)
        dispatch({ type: "SET_WORKFLOWS", data: result });
      } else {
        dispatch({ type: "SET_ERROR", data: res });
      }
    });
  }

  function onSearchInputChange(value: string) {
    setSearchValue(value);
  }

  function onSearchButtonClick() {
    const filtered =
      searchValue === ''
        ? repositories
        : repositories.filter(item => {
          const str = (typeof item.repository_name === 'string') ? item.repository_name : item.repository_name;
          return str.toLowerCase().indexOf(searchValue.toLowerCase()) !== -1;
        });

    setFilteredItems(filtered || []);
  }

  if (typeof repositories[0] !== 'undefined') {
    firstItemText = repositories[0].repository_name
  }

  // Validates if the repository and organization are correct
  const validParams = (repository, organization) => {
    if (isValidTeam()) {
      if (repositories.find(r => r.git_organization == organization && r.repository_name == repository)) {
        return true;
      }
      if (repository == "default" && organization == "") {
        return true;
      }
    }
    return false;
  }

  return (
    <PageSection>
      <div style={{ backgroundColor: "white", paddingTop: "5px" }}>
        <div style={{ padding: "10px" }}>
          <ContextSelector
            toggleText={selected}
            onSearchInputChange={onSearchInputChange}
            isOpen={isOpen}
            searchInputValue={searchValue}
            onToggle={onToggle}
            onSearchButtonClick={onSearchButtonClick}
            screenReaderLabel="Selected Project:"
            isPlain
            isText
          >
            {repositories.map((item, index) => {
              const [text] = (typeof item === 'string')
                ? [item]
                : [item];
              return <ContextSelectorItem key={index} onClick={() => {
                setSelected(text.repository_name);
                setOrg(text.git_organization)
                setOpen(!isOpen);
                getWorkflows(text.repository_name);
                params.set("repository", text.repository_name)
                params.set("organization", text.git_organization)
                history.push(window.location.pathname + '?' + params.toString());
              }} >{text.repository_name}</ContextSelectorItem>;
            })}
          </ContextSelector>
          <Button onClick={() => navigator.clipboard.writeText(window.location.href)} variant="link" icon={<CopyIcon />} iconPosition="right">
            Copy link
          </Button>
        </div>
        <hr />
        {validParams(selected, org) && <TableComposable aria-label="Actions table" style={{ padding: "10px" }}>
          <Caption>All GitHub Actions available in the repository {selected}</Caption>
          <Thead>
            <Tr>
              <Th>{columnNames.name}</Th>
              <Th>{columnNames.job}</Th>
              <Th>{columnNames.badge}</Th>
              <Th>{columnNames.html_url}</Th>
            </Tr>
          </Thead>
          <Tbody>
            {workflows.map(workf => {
              const workflow: Workflows = workf
              return (
                <Tr key={workflow.workflow_name || ""}>
                  <Td dataLabel={columnNames.name}>{workflow.workflow_name || ""}</Td>
                  <Td dataLabel={columnNames.job}>{workflow.job || ""}</Td>
                  <Td dataLabel={columnNames.badge}><img src={workflow.badge_url}></img></Td>
                  <Td dataLabel={columnNames.html_url}><a href={workflow.html_url}><ExternalLinkAltIcon>Link</ExternalLinkAltIcon>Go to job</a></Td>
                </Tr>
              );
            })}
          </Tbody>
        </TableComposable>
        }
        {!validParams(selected, org) && <EmptyState variant={EmptyStateVariant.xl}>
          <EmptyStateIcon icon={ExclamationCircleIcon} />
          <Title headingLevel="h1" size="lg">
            Something went wrong. Please, check the URL.
          </Title>
        </EmptyState>
        }
      </div>
    </PageSection>

  );
};