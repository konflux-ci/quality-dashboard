import React from 'react';
import {
  Title,
  EmptyState,
  EmptyStateIcon,
  EmptyStateBody,
} from '@patternfly/react-core';
import CubesIcon from '@patternfly/react-icons/dist/esm/icons/cubes-icon';

export const Empty: React.FunctionComponent<{title: string, body: string}> = ({title, body}) => (
  <EmptyState>
    <EmptyStateIcon icon={CubesIcon} />
    <Title headingLevel="h4" size="lg">
      {title}
    </Title>
    <EmptyStateBody>
      {body}
    </EmptyStateBody>
  </EmptyState>
);
