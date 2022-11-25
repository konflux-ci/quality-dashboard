import React, { } from 'react';
import { TableComponent, TableComponentProps } from '@app/Repositories/TableComponent';
import { FormModal, ModalContext, useDefaultModalContextState } from '@app/Repositories/CreateRepository';

type RepositoriesTableProps = TableComponentProps

export const RepositoriesTable = ({showCoverage, showDiscription, showTableToolbar, enableFiltersOnTheseColumns} : RepositoriesTableProps) => {

    const defaultModalContext = useDefaultModalContextState();

    return (
        <ModalContext.Provider value={defaultModalContext}>
            <React.Fragment>
                <TableComponent showCoverage={showCoverage} showDiscription={showDiscription} showTableToolbar={showTableToolbar} enableFiltersOnTheseColumns={enableFiltersOnTheseColumns}/>
                <FormModal></FormModal>
            </React.Fragment>
        </ModalContext.Provider>
    );
};
