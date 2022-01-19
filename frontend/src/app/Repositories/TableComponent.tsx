/* eslint-disable no-console */
import React, { useContext, useEffect, useState } from 'react';
import {
  TableComposable,
  Thead,
  Tr,
  Th,
  Tbody,
  Td,
  ActionsColumn,
  IAction,
  Caption
} from '@patternfly/react-table';
import { Context } from '@app/store/store';
import { deleteRepositoryAPI, getRepositories } from '@app/utils/APIService';
import { ExternalLinkAltIcon } from '@patternfly/react-icons';

interface Repository {
  git_organization: string;
  repository_name: string;
  git_url: string
  description: string
}

const columnNames = {
  organization: 'GitHub Organization',
  repository: 'Repository',
  description: 'Description',
};

const defaultActions = (repo: Repository): IAction[] => [
  {
    title: 'Delete Repository',
    onClick: () => deleteRepository(repo.git_organization, repo.repository_name)
  },
];

async function deleteRepository(gitOrg:string, repoName:string) {
  const data = {
    git_organization: gitOrg,
    repository_name : repoName,
  }
  try {
    await deleteRepositoryAPI(data)

    window.location.reload();
  } catch (error) {
    console.log(error)
  }
}

export const TableComponent: React.FunctionComponent = () => {
  const [repos, setRepositories] = useState([])
  const { state, dispatch } = useContext(Context) 

  useEffect(()=> {
    getRepositories().then((res) => {
      if(res.code === 200) {
          const result = res.data;
          dispatch({ type: "SET_REPOSITORIES", data: result });
      } else {
          dispatch({ type: "SET_ERROR", data: res });
      }
    });
  }, [repos, setRepositories, dispatch])

  const repositories: Repository[] = state.repositories

  return (
    <React.Fragment>
      <TableComposable aria-label="Actions table">
      <Caption>Repositories Summary</Caption>
        <Thead>
          <Tr>
            <Th>{columnNames.organization}</Th>
            <Th>{columnNames.repository}</Th>
            <Th>{columnNames.description}</Th>
          </Tr>
        </Thead>
        <Tbody>
          {repositories.map(repo => {
            const rowActions: IAction[] | null = defaultActions(repo);
            return (
              <Tr key={repo.repository_name}>
                <Td dataLabel={columnNames.organization}>{repo.git_organization}</Td>
                <Td dataLabel={columnNames.repository}><a href= {repo.git_url}>{repo.repository_name}</a><a href={repo.git_url}><ExternalLinkAltIcon style={{marginLeft: "1%"}}></ExternalLinkAltIcon></a></Td>
                <Td dataLabel={columnNames.description}>{repo.description}</Td>
                <Td>
                  {rowActions ? (
                    <ActionsColumn
                      items={rowActions}
                    />
                  ) : null}
                </Td>
              </Tr>
            );
          })}
        </Tbody>
      </TableComposable>
    </React.Fragment>
  );
};
