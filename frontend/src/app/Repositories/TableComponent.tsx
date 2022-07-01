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
  Caption,
  ThProps
} from '@patternfly/react-table';
import { 
  Alert, Pagination, PaginationVariant, 
  Button, ButtonVariant,
  Toolbar, ToolbarItem, ToolbarContent,
} from '@patternfly/react-core';
import { Context } from '@app/store/store';
import { deleteRepositoryAPI, getRepositories } from '@app/utils/APIService';
import { ExternalLinkAltIcon } from '@patternfly/react-icons';
import _ from 'lodash';
import { useModalContext } from './CreateRepository';
import { Repository } from './Repositories';

export interface TableComponentProps {
  showCoverage?: boolean
  showDiscription?: boolean
  showTableToolbar?: boolean
}

const columnNames = {
  organization: 'GitHub Organization',
  repository: 'Repository',
  description: 'Description',
  coverageType: 'Coverage Type',
  coverage: 'Code Covered',
};

const rederCoverageEffects = (repo: Repository) => {
  const coveredFixed = repo.code_coverage.coverage_percentage
  if (coveredFixed >= 0 && coveredFixed <= 33.33 ) {
    return <Alert title={coveredFixed.toFixed(2)+"%"} variant="danger" isInline isPlain />
  } else if (coveredFixed >= 33.33 && coveredFixed <= 66.66) {
    return <Alert title={coveredFixed.toFixed(2)+"%"} variant="warning" isInline isPlain />
  }
  return <Alert title={coveredFixed.toFixed(2)+"%"} variant="success" isInline isPlain />
}

export const TableComponent = ({showCoverage, showDiscription, showTableToolbar}: TableComponentProps) => {
  const { state, dispatch } = useContext(Context) 
  const [perpage, onperpageset] = useState(10)
  const [repos, setRepositories] = useState<any>([])
  const [page, onPageset]= useState(1)
  const [allreps, setallreps] = useState<any>(state.Allrepositories);
  const modalContext = useModalContext()

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
  
  async function editRepository(repo: Repository) {
    try {
      modalContext.handleModalToggle(true, repo)
    } catch (error) {
      console.log(error)
    }
  }

  const defaultActions = (repo: Repository): IAction[] => [
    {
      title: 'Delete Repository',
      onClick: () => deleteRepository(repo.git_organization, repo.repository_name)
    },
    {
      title: 'Edit Repository',
      onClick: () => editRepository(repo)
    },
  ];

  function onPageselect(e, page){
    onPageset(page)
  }

  function onperpageselect(e, Perpage){
    onperpageset(Perpage)
  }

  // Sort helpers
  const [activeSortIndex, setActiveSortIndex] = React.useState<number | null>(null);
  const [activeSortDirection, setActiveSortDirection] = React.useState<'asc' | 'desc' | null>(null);
  const getSortableRowValues = (repo: Repository): (string | number)[] => {
    const { git_organization, repository_name, code_coverage } = repo;
    return [git_organization, repository_name, code_coverage.coverage_percentage];
  };
  let sortedRepositories = repos;
  if (activeSortIndex !== null) {
    sortedRepositories = repos.sort((a, b) => {
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

  useEffect(()=> {
    getRepositories(perpage).then((res)=> {
      if(res.code === 200) {
        const result = res.data;
        setallreps(res.all)
        let repositories: any = result
        const repos = repositories[page-1]
        const ress : any = []
        _.each(repos, function(ele, index, array){
          ress.push(ele)
        })
        setRepositories(ress)
        onPageset(page)
        dispatch({ type: "SET_REPOSITORIES", data: result });
        dispatch({type: "SET_REPOSITORIES_ALL", data: res.all});
      } else {
          dispatch({ type: "SET_ERROR", data: res });
      }
    })
  }, [page, perpage, setRepositories, dispatch])

  return (
    <React.Fragment>
      {showTableToolbar && 
        <Toolbar style={{marginBottom: "5px"}}>
          <ToolbarContent>
            <ToolbarItem>
              <Button variant={ButtonVariant.secondary} onClick={modalContext.handleModalToggle}>
                Add Git Repository
              </Button>
            </ToolbarItem>
          </ToolbarContent>
        </Toolbar>
      }
      <TableComposable aria-label="Actions table">
      <Caption>Repositories Summary</Caption>
        <Thead>
          <Tr>
          <Th sort={getSortParams(0)}>{columnNames.organization}</Th>
           <Th sort={getSortParams(1)}>{columnNames.repository}</Th>
           {showDiscription && 
           <Th>{columnNames.description}</Th>
           }
           {showCoverage &&
            <Th>{columnNames.coverageType}</Th>
           }
           {showCoverage &&
            <Th sort={getSortParams(2)}>{columnNames.coverage}</Th>
           }  
          </Tr>
        </Thead>
        <Tbody>
          {repos.map(repo => {
            const rowActions: IAction[] | null = defaultActions(repo);
            return (
              <Tr key={repo.repository_name}>
                <Td>{repo.git_organization}</Td>
                <Td><a href= {repo.git_url}>{repo.repository_name}</a><a href={repo.git_url}><ExternalLinkAltIcon style={{marginLeft: "1%"}}></ExternalLinkAltIcon></a></Td>
                {showDiscription && 
                <Td>{repo.description}</Td>
                }
                {showCoverage && 
                  <Td><a href={`https://app.codecov.io/gh/${repo.git_organization}/${repo.repository_name}`}>CodeCov<ExternalLinkAltIcon style={{marginLeft: "0.5%"}}></ExternalLinkAltIcon></a></Td>
                }
                {showCoverage && 
                  <Td>{rederCoverageEffects(repo)}</Td>
                }
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
      <Pagination
        itemCount={allreps.length}
        perPage={perpage}
        page={page}
        variant={PaginationVariant.bottom}
        onSetPage={onPageselect}
        onPerPageSelect={onperpageselect}
        
      />
    </React.Fragment>
  );
};
