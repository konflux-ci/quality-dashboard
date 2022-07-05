import React, { } from 'react';
import { PageSection } from '@patternfly/react-core';
import { TableComponent, TableComponentProps } from './TableComponent';
import { FormModal, ModalContext, useDefaultModalContextState } from './CreateRepository';

interface RepositoriesTableProps extends TableComponentProps {}

export const RepositoriesTable = ({showCoverage, showDiscription, showTableToolbar, enableFiltersOnTheseColumns} : RepositoriesTableProps) => {

    const defaultModalContext = useDefaultModalContextState()

    return (
        <ModalContext.Provider value={defaultModalContext}>
            <React.Fragment>
                <TableComponent showCoverage={showCoverage} showDiscription={showDiscription} showTableToolbar={showTableToolbar} enableFiltersOnTheseColumns={enableFiltersOnTheseColumns}/>
                <FormModal></FormModal>
            </React.Fragment>
        </ModalContext.Provider>
    );
};
