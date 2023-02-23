import React, { useState, useContext, useEffect } from 'react';
import { createTeam, createRepository, getTeams } from "@app/utils/APIService";
import {
    Wizard,
    TextInput,
    FormGroup,
    Form,
    TextArea,
    DescriptionList,
    DescriptionListGroup,
    DescriptionListDescription,
    DescriptionListTerm,
    Title,
    Spinner,
    Alert,
    AlertGroup,
    AlertVariant,
} from '@patternfly/react-core';
import { ReactReduxContext } from 'react-redux';

interface AlertInfo {
    title: string;
    variant: AlertVariant;
    key: string;
}


export const TeamForm = (props) => {
    const { store } = useContext(ReactReduxContext);
    const state = store.getState();
    const dispatch = store.dispatch;

    const [stepIdReached, setState] = useState<string>("team");
    const [newTeamName, setNewTeamName] = useState<string>("");
    const [newTeamDesc, setNewTeamDesc] = useState<string>("");
    const [newRepoName, setNewRepoName] = useState<string>("");
    const [newOrgName, setNewOrgName] = useState<string>("");
    const [creationLoading, setCreationLoading] = useState<boolean>(false);
    const [alerts, setAlerts] = React.useState<AlertInfo[]>([]);
    const [creationError, setCreationError] = useState<boolean>(false);
    const [isFinishedWizard, setIsFinishedWizard] = useState<boolean>(false);
    const [isOpen, setOpen] = useState<boolean>(false);

    useEffect(
        () => {
            setOpen(props.isOpen);
        },
        [props.isOpen],
    );


    const onSubmit = async () => {
        // Create a team
        setIsFinishedWizard(true)
        setCreationLoading(true)
        const data = {
            "team_name": newTeamName,
            "description": newTeamDesc,
        }

        createTeam(data).then(response => {
            if (response.code == 200) {
                setAlerts(prevAlertInfo => [...prevAlertInfo, {
                    title: 'Team created',
                    variant: AlertVariant.success,
                    key: "team-created"
                }]);
            } else {
                setAlerts(prevAlertInfo => [...prevAlertInfo, {
                    title: 'Could not create new team',
                    variant: AlertVariant.danger,
                    key: "team-not-created"
                }]);
                setCreationError(true)
            }
        }).catch(error => {
            console.log("NEW TEAM", data);
            console.log(error);
        })


        // Create a repo optionally (only if fields are populated)
        if (newRepoName != "" && newOrgName != "") {
            try {
                const data = {
                    git_organization: newOrgName,
                    repository_name: newRepoName,
                    jobs: {
                        github_actions: {
                            monitor: false
                        }
                    },
                    artifacts: [],
                    team_name: newTeamName
                }
                const response = await createRepository(data)
                if (response.code == 200) {
                    setAlerts(prevAlertInfo => [...prevAlertInfo, {
                        title: 'Repository created',
                        variant: AlertVariant.success,
                        key: "repo-created"
                    }]);
                } else {
                    setAlerts(prevAlertInfo => [...prevAlertInfo, {
                        title: 'Could not add repository',
                        variant: AlertVariant.danger,
                        key: "repo-not-created"
                    }]);
                    setCreationError(true)
                }
            }
            catch (error) {
                console.log(error)
            }
        }

        setCreationLoading(false)
        if (!creationError) {
            setAlerts(prevAlertInfo => [...prevAlertInfo, {
                title: 'You resources have created successfully. You can close the modal now.',
                variant: AlertVariant.info,
                key: "all-created"
            }]);
            getTeams().then(data => {
                if (data.data.length > 0) {
                    dispatch({ type: "SET_TEAM", data: newTeamName });
                    dispatch({ type: "SET_TEAMS_AVAILABLE", data: data.data });
                }
            })
        }

    };

    const onNext = (id) => {
        setState(id.id);
    };

    const onBack = (id) => {
        setState(id.id);
    };

    function onClear() {
        setOpen(false)
        props.handleChange(event, false);
        setNewTeamName("")
        setNewTeamDesc("")
        setNewOrgName("")
        setNewRepoName("")
    };

    const TeamData = (
        <div className={'pf-u-m-lg'} >
            <Form>
                <FormGroup label="Team Name" isRequired fieldId="team-name" helperText="Include the name for your team">
                    <TextInput value={newTeamName} type="text" onChange={(value) => { setNewTeamName(value) }} aria-label="text input example" placeholder="Include the name for your team" />
                </FormGroup>
                <FormGroup label="Description" isRequired fieldId='team-description' helperText="Include a description for your team">
                    <TextArea value={newTeamDesc} type="text" onChange={(value) => { setNewTeamDesc(value) }} aria-label="text area example" placeholder="Include a description for your team" />
                </FormGroup>
            </Form>
        </div>
    )

    const AddRepo = (
        <div className={'pf-u-m-lg'} >
            <Title headingLevel="h6" size="xl">Optionally: add a repository to your team</Title>
            <Form>
                <FormGroup label="Repository Name" fieldId="repo-name" helperText="Add a repository">
                    <TextInput value={newRepoName} type="text" onChange={value => setNewRepoName(value)} aria-label="text input example" placeholder="Add a repository" />
                </FormGroup>
                <FormGroup label="Organization Name" fieldId="org-name" helperText="Specify the organization">
                    <TextInput value={newOrgName} type="text" onChange={value => setNewOrgName(value)} aria-label="text input example" placeholder="Specify the organization" />
                </FormGroup>
            </Form>
        </div>
    )

    const DataReview = (
        <div>
            <Title headingLevel="h6" size="xl">Review your data</Title>
            <div style={{ marginTop: '2em' }}>
                <DescriptionList isHorizontal>
                    <DescriptionListGroup>
                        <DescriptionListTerm>Team Name</DescriptionListTerm>
                        <DescriptionListDescription>{newTeamName}</DescriptionListDescription>
                    </DescriptionListGroup>
                    <DescriptionListGroup>
                        <DescriptionListTerm>Team Description</DescriptionListTerm>
                        <DescriptionListDescription>{newTeamDesc}</DescriptionListDescription>
                    </DescriptionListGroup>
                    <DescriptionListGroup>
                        <DescriptionListTerm>Repository</DescriptionListTerm>
                        <DescriptionListDescription>{newRepoName}</DescriptionListDescription>
                    </DescriptionListGroup>
                    <DescriptionListGroup>
                        <DescriptionListTerm>Organization</DescriptionListTerm>
                        <DescriptionListDescription>{newOrgName}</DescriptionListDescription>
                    </DescriptionListGroup>
                </DescriptionList>
            </div>
            <div style={{ marginTop: "2em" }}>
                {creationLoading && <Spinner isSVG aria-label="Contents of the basic example" />}
                <AlertGroup isLiveRegion aria-live="polite" aria-relevant="additions text" aria-atomic="false">
                    {alerts.map(({ title, variant, key }) => (
                        <Alert variant={variant} isInline isPlain title={title} key={key} />
                    ))}
                </AlertGroup>
            </div>
        </div>
    )

    const ValidateTeam = () => { return newTeamName != "" && newTeamDesc != "" }
    const ValidateRepoAndOrg = () => {
        if (newRepoName != "" || newOrgName != "") {
            return newRepoName != "" && newOrgName != ""
        }
        else if (newRepoName == "" && newOrgName == "") {
            return true
        }
        return false

    }

    const steps = [
        { id: 'team', name: 'Team Name', component: TeamData, enableNext: ValidateTeam() },
        { id: 'repo', name: 'Add a repository', component: AddRepo, canJumpTo: ValidateTeam(), enableNext: ValidateRepoAndOrg() },
        {
            id: 'review',
            name: 'Review',
            component: DataReview,
            nextButtonText: 'Create',
            canJumpTo: (ValidateTeam() && ValidateRepoAndOrg()),
            hideCancelButton: true,
            isFinishedStep: isFinishedWizard
        }
    ];

    const title = 'Create new Team';

    return (
        <React.Fragment>
            <Wizard
                steps={steps}
                onNext={onNext}
                onBack={onBack}
                onClose={onClear}
                onSave={onSubmit}
                cancelButtonText="Close"
                height={600}
                title={title}
                description="Create a new team and (optionally) add a new repository"
                isOpen={isOpen}
                hideClose={false}
            />
        </React.Fragment>
    )
}
