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
import { Alert, Pagination, PaginationVariant } from '@patternfly/react-core';
import { Context } from '@app/store/store';
import { deleteRepositoryAPI, getRepositories } from '@app/utils/APIService';
import { ExternalLinkAltIcon } from '@patternfly/react-icons';
import _ from 'lodash';
import { element } from 'prop-types';

type TableComponentProps = {
  showCoverage: boolean
  showDiscription: boolean
}
interface Coverage {
  coverage_percentage: number
}
interface Repository {
  git_organization: string;
  repository_name: string;
  git_url: string;
  description: string;
  code_coverage: Coverage;
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

export const TableComponent = ({showCoverage, showDiscription}: TableComponentProps) => {
  const { state, dispatch } = useContext(Context) 
  const [perpage, onperpageset] = useState(10)
  const [repos, setRepositories] = useState<any>([])
  const [page, onPageset]= useState(1)
  const [allreps, setallreps] = useState<any>(state.Allrepositories);
  
  function onPageselect(e, page){
    onPageset(page)
  }
  function onperpageselect(e, Perpage){
    onperpageset(Perpage)
  }

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
      <TableComposable aria-label="Actions table">
      <Caption>Repositories Summary</Caption>
        <Thead>
          <Tr>
          <Th>{columnNames.organization}</Th>
           <Th>{columnNames.repository}</Th>
           {showDiscription && 
           <Th>{columnNames.description}</Th>
           }
           {showCoverage &&
            <Th>{columnNames.coverageType}</Th>
           }
           {showCoverage &&
            <Th>{columnNames.coverage}</Th>
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
