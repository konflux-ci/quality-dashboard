import React, { } from 'react';
import { TableComponent } from './TableComponent';

export const TeamsTable: React.FunctionComponent = () => {

  //  const defaultModalContext = useDefaultModalContextState();

    return (
        // <ModalContext.Provider value={defaultModalContext}>
            <React.Fragment>
                <TableComponent/>
                {/* <FormModal></FormModal> */}
            </React.Fragment>
        // </ModalContext.Provider>
    );
};
