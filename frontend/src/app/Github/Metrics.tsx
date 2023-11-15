import React, { SetStateAction, useContext } from 'react';
import { Grid, GridItem, Modal } from '@patternfly/react-core';
import { GitHubInfoCard } from '@app/utils/sharedComponents';
import { CodeCov } from './CodeCov';
import { GithubActions } from './GithubActions';
import { PullRequestCard, PullRequestsGraphic } from './PullRequests';


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
    const repo = modalContext.data.value;
    const prs = repo?.prs
    const workflows = repo?.workflows
    const repoName = repo?.repository_name
    const repoOrg = repo?.git_organization
    
    return (
        <React.Fragment>
            <Modal
                style={{ backgroundColor: '#D1D1D1' }}
                title={'Metrics'}
                isOpen={modalContext.isModalOpen.value}
                onClose={modalContext.handleModalToggle}

            >
                <Grid hasGutter>
                    <GridItem span={3} rowSpan={2}>
                        <GitHubInfoCard data={[{ title: 'Repository', value: repoName }, { title: 'Organization', value: repoOrg }, { title: 'Description', value: repo.description },]} org={repoOrg} repoName={repoName}></GitHubInfoCard>
                    </GridItem>

                    <GridItem span={7} rowSpan={4}>
                        <PullRequestsGraphic metrics={prs?.metrics} summary={prs?.summary}></PullRequestsGraphic>
                    </GridItem>

                    <GridItem span={2} rowSpan={1}>
                        <PullRequestCard title="Time To Merge PR Avg Days" subtitle="Selected Time Range" total={prs?.summary?.merge_avg}></PullRequestCard>
                    </GridItem>

                    <GridItem span={2} rowSpan={1}>
                        <PullRequestCard
                            title="Retest Avg in Open PRs"
                            subtitle="Selected Time Range"
                            total={prs?.summary?.retest_avg}
                        ></PullRequestCard>
                    </GridItem>

                    <GridItem span={2} rowSpan={1}>
                        <PullRequestCard title="Open PRs" subtitle="Total" total={prs?.summary?.open_prs}></PullRequestCard>
                    </GridItem>

                    <GridItem span={1} rowSpan={1}>
                        <PullRequestCard title="Created PRs" subtitle="Selected Time Range" total={prs?.summary?.created_prs_in_time_range}></PullRequestCard>
                    </GridItem>

                    <GridItem span={2} rowSpan={1}>
                        <PullRequestCard
                            title="Retest Before Merge Avg"
                            subtitle="Selected Time Range"
                            // edge case of service-provider-integration-operator
                            // https://github.com/redhat-appstudio/service-provider-integration-operator/pull/548#issuecomment-1494149514
                            total={prs?.summary?.retest_before_merge_avg == 0.01 ? 'N/A' : prs?.summary?.retest_before_merge_avg}
                        ></PullRequestCard>
                    </GridItem>

                    <GridItem span={2} rowSpan={1}>
                        <PullRequestCard title="Merged PRs" subtitle="Total" total={prs?.summary?.merged_prs}></PullRequestCard>
                    </GridItem>
                    
                    <GridItem span={1} rowSpan={1}>
                        <PullRequestCard title="Merged PRs" subtitle="Selected Time Range" total={prs?.summary?.merged_prs_in_time_range}></PullRequestCard>
                    </GridItem>

                    <GridItem span={2} rowSpan={1}>
                        <CodeCov repo={repo}></CodeCov>
                    </GridItem>

                    {workflows?.length > 0 && (
                        <GridItem span={12}>
                            <GithubActions repoName={repoName} workflows={workflows}></GithubActions>
                        </GridItem>
                    )}
                </Grid>
            </Modal>
        </React.Fragment>
    );
}