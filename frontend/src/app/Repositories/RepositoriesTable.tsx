import React, { } from 'react';
import { PageSection } from '@patternfly/react-core';
import { TableComponent, TableComponentProps } from './TableComponent';
import { FormModal, ModalContext, useDefaultModalContextState } from './CreateRepository';

interface RepositoriesTableProps extends TableComponentProps {}

export const RepositoriesTable = ({showCoverage, showDiscription, showTableToolbar} : RepositoriesTableProps) => {

    const defaultModalContext = useDefaultModalContextState()

    return (
        <ModalContext.Provider value={defaultModalContext}>
            <React.Fragment>
                <PageSection>
                    <TableComponent showCoverage={showCoverage} showDiscription={showDiscription} showTableToolbar={showTableToolbar}/>
                    <FormModal></FormModal>
                </PageSection>
            </React.Fragment>
        </ModalContext.Provider>
    );
};
