import React, { } from 'react';
import { TableComponent, TableComponentProps } from './TableComponent';
import { FormModal, ModalContext, useDefaultModalContextState } from './CreateRepository';

type RepositoriesTableProps = TableComponentProps

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
