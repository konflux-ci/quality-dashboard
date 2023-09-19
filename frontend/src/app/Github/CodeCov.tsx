import React from 'react';
import { Card, CardBody, CardTitle } from '@patternfly/react-core';
import { isUndefined } from 'lodash';
import { EmptyStateIcon } from '@patternfly/react-core';
import { TrendDownIcon, TrendUpIcon } from '@patternfly/react-icons';
import { help } from './PullRequests';


export interface Coverage {
  coverage_percentage: number;
  coverage_trend: string;
}

const getColor = (coveredFixed) => {
  if (coveredFixed >= 0 && coveredFixed <= 33.33) {
    return 'red';
  } else if (coveredFixed >= 33.33 && coveredFixed <= 66.66) {
    return 'orange';
  }
  return 'green';
};

export const CodeCov = (props) => {
  const repo = props.repo;

  return (
    <Card style={{ width: '100%', height: '100%', textAlign: 'center' }}>
      <CardTitle>
        <div>CodeCov</div>
        <div style={{ color: 'grey', fontSize: 12 }}>
          {"Two last commits"}
          {help("The coverage trend is calculated through the two last commits. No trend arrow means that the coverage trend is stable")}
        </div>
      </CardTitle>
      <CardBody>
        <div style={{ fontSize: 25 }}>
          {!isUndefined(repo) && (repo?.coverage?.coverage_percentage != 0) ? GetCodeCovInfo(repo, 'center'): 'N/A'}
        </div>
      </CardBody>
    </Card>
  );
};

export const GetCodeCovInfo = (repo, justifyContent) => {
  const coverage = repo.code_coverage.coverage_percentage
  const coverageTrend = repo.coverage_trend

  return (
    <div>
      <div style={{ display: "flex", justifyContent: justifyContent }}>
       
       <div style={{ color: getColor(coverage) }}>{coverage + '%'}</div>

        {coverageTrend == "descending" && <div style={{ marginLeft: 5 }}><EmptyStateIcon icon={TrendDownIcon} /></div>}
        {coverageTrend == "ascending" && <div style={{ marginLeft: 5 }}><EmptyStateIcon icon={TrendUpIcon} /></div>}

      </div>
      <div style={{ color: 'grey', fontSize: 14 }}>
      <a
        href={`https://app.codecov.io/gh/${repo.git_organization}/${repo.repository_name}/commits`}
        target="blank"
        rel="noopener noreferrer"
      >
        More info
      </a>
    </div>
    </div >
  );
};