import React from 'react';
import { Card, CardBody, CardTitle } from '@patternfly/react-core';
import { isUndefined } from 'lodash';

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
          {!isUndefined(repo) ? (
            <a
              href={`https://app.codecov.io/gh/${repo.organization}/${repo.repoName}`}
              target="blank"
              rel="noopener noreferrer"
            >
              More info
            </a>
          ) : (
            'N/A'
          )}
        </div>
      </CardTitle>
      <CardBody>
        <div style={{ fontSize: 25 }}>
          {!isUndefined(repo) ? <div style={{ color: getColor(repo.coverage.coverage_percentage) }}>{repo.coverage.coverage_percentage.toFixed(2) + '%'}</div> : 'N/A'}
        </div>
      </CardBody>
    </Card>
  );
};
