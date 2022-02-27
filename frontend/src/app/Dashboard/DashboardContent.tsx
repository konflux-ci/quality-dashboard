import React, { useContext, useEffect, useState } from 'react';
import { TableComposable, Thead, Tbody, Tr, Th, Td, Caption } from '@patternfly/react-table';
import { Context } from '@app/store/store';
import { getRepositories } from '@app/utils/APIService';
import { Alert, PageSection, Pagination, PaginationVariant } from '@patternfly/react-core';
import { ExternalLinkAltIcon, ArrowCircleDownIcon } from '@patternfly/react-icons';
import { number } from 'prop-types';

interface Coverage {
  coverage_percentage: number
}
interface Repository {
  git_organization: string;
  repository_name: string;
  git_url: string
  code_coverage: Coverage
}

const columnNames = {
  organization: 'GitHub Organization',
  repository: 'Repository',
  coverageType: 'Coverage Type',
  coverage: 'Code Covered',
};

const rederCoverageEffects = (repo: Repository) => {
  const coveredFixed = repo.code_coverage.coverage_percentage
  if (coveredFixed >= 0 && coveredFixed <= 33.33 ) {
    return <Td dataLabel={columnNames.coverage}><Alert title={coveredFixed.toFixed(2)+"%"} variant="danger" isInline isPlain /></Td>
  } else if (coveredFixed >= 33.33 && coveredFixed <= 66.66) {
    return <Td dataLabel={columnNames.coverage}><Alert title={coveredFixed.toFixed(2)+"%"} variant="warning" isInline isPlain /></Td>
  }
  return <Td dataLabel={columnNames.coverage}><Alert title={coveredFixed.toFixed(2)+"%"} variant="success" isInline isPlain /></Td>
}

export const DashboardContent = () => {
  const [ gitRepos, setRepositories ] = useState([])
  const { state, dispatch } = useContext(Context)

  useEffect(()=> {
    getRepositories().then((res)=> {
      if(res.code === 200) {
        console.log(res.data)
        const result = res.data;
        dispatch({ type: "SET_REPOSITORIES", data: result });
      } else {
          dispatch({ type: "SET_ERROR", data: res });
      }
    })
  }, [gitRepos, setRepositories, dispatch])

  const repositories: Repository[] = state.repositories

  return (
    <React.Fragment>
      <PageSection style={{
        minHeight : "12%"
      }}>
        <TableComposable aria-label="Actions table">
        <Caption>Repositories Summary</Caption>
          <Thead>
            <Tr>
              <Th>{columnNames.organization}</Th>
              <Th>{columnNames.repository}</Th>
              <Th>{columnNames.coverageType}</Th>
              <Th>{columnNames.coverage}</Th>
            </Tr>
          </Thead>
          <Tbody>
            {repositories.map(repo => {
              return (
                <Tr key={repo.repository_name}>
                  <Td dataLabel={columnNames.organization}>{repo.git_organization}</Td>
                  <Td dataLabel={columnNames.repository}><a href= {repo.git_url}>{repo.repository_name}</a><a href={repo.git_url}><ExternalLinkAltIcon style={{marginLeft: "0.5%"}}></ExternalLinkAltIcon></a></Td>
                  <Td dataLabel={columnNames.coverageType}><a href={`https://app.codecov.io/gh/${repo.git_organization}/${repo.repository_name}`}>CodeCov<ExternalLinkAltIcon style={{marginLeft: "0.5%"}}></ExternalLinkAltIcon></a></Td>
                  {rederCoverageEffects(repo)}
                </Tr>
              );
            })}
          </Tbody>
        </TableComposable>
      </PageSection>
    </React.Fragment>
  );
}
