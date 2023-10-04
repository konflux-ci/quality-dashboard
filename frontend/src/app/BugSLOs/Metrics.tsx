import React, { SetStateAction, useContext } from 'react';
import { Card, CardBody, CardTitle, Grid, GridItem, Modal } from '@patternfly/react-core';
import { BugSLOTable } from './BugSLOsTable';


interface IModalContext {
    isModalOpen: IModalContextMember;
    handleModalToggle;
    data: IModalContextMember;
}

interface IModalContextMember {
    value: any;
    set: SetStateAction<any>;
}

export const MetricsModalContext = React.createContext<IModalContext>({
    isModalOpen: { set: undefined, value: false },
    handleModalToggle: () => { },
    data: { set: undefined, value: false },
});

export const useMetricsModalContext = () => {
    return useContext(MetricsModalContext);
};

export const useMetricsModalContextState = () => {
    const [isModalOpen, setModalOpen] = React.useState(false);
    const [data, setData] = React.useState({});
    const defaultModalContext = useMetricsModalContext();

    defaultModalContext.isModalOpen = { set: setModalOpen, value: isModalOpen };
    defaultModalContext.data = { set: setData, value: data };
    defaultModalContext.handleModalToggle = (data: any) => {
        defaultModalContext.isModalOpen.set(!defaultModalContext.isModalOpen.value);
        defaultModalContext.data.set(data);
    };
    return defaultModalContext;
};

export const GetMetrics = () => {
    const modalContext = useMetricsModalContext();
    const project = modalContext.data.value;
    const bugSLOs = project?.bug_slos

    const CustomCard = (props) => {
        const title = props.title
        const body = props.body

        return (
            <Card style={{ width: '100%', height: '100%', textAlign: 'center' }}>
                <CardTitle>{title}</CardTitle>
                <CardBody>{body}</CardBody>
            </Card>
        )
    }

    return (
        <React.Fragment>
            <Modal
                style={{ backgroundColor: '#D1D1D1' }}
                title={'Metrics'}
                isOpen={modalContext.isModalOpen.value}
                onClose={modalContext.handleModalToggle}
            >
                <Grid hasGutter>
                    <GridItem span={4} rowSpan={4}>
                        <CustomCard
                            title="Days Without Priority Avg"
                            body={project?.red_triage_time_bug_slo_info?.average + " days"}
                        >
                        </CustomCard>
                    </GridItem>
                    <GridItem span={4} rowSpan={4}>
                        <CustomCard
                            title="Days Without Assignee Avg"
                            body={project?.red_response_time_bug_slo_info?.average + " days"}
                        >
                        </CustomCard>
                    </GridItem>
                    <GridItem span={4} rowSpan={4}>
                        <CustomCard
                            title="Days Without Resolution Avg"
                            body={project?.red_resolution_time_bug_slo_info?.average + " days"}
                        >
                        </CustomCard>
                    </GridItem>
                    <GridItem span={12}>
                        <BugSLOTable bugSLOs={bugSLOs}></BugSLOTable>
                    </GridItem>
                </Grid>
            </Modal>
        </React.Fragment>
    );
}