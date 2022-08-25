import React, {useState} from 'react';
import { createTeam, createRepository } from "@app/utils/APIService";
import { 
  Wizard, PageSection, PageSectionVariants, 
  TextContent, Text, TextVariants,TextInput, FormGroup, Form, TextArea,
  DescriptionList, DescriptionListGroup, DescriptionListDescription, DescriptionListTerm, Title, Spinner,
  Alert, AlertGroup, AlertVariant,
} from '@patternfly/react-core';

import { useHistory } from 'react-router-dom';

interface AlertInfo {
  title: string;
  variant: AlertVariant;
  key: string;
}

export const TeamsWizard = () => {
  const history = useHistory()  

  const [ stepIdReached, setState ] = useState<string>("team");
  const [ newTeamName, setNewTeamName ] = useState<string>("");
  const [ newTeamDesc, setNewTeamDesc ] = useState<string>("");
  const [ newRepoName, setNewRepoName ] = useState<string>("");
  const [ newOrgName, setNewOrgName ] = useState<string>("");
  const [ creationLoading, setCreationLoading ] = useState<boolean>(false);
  const [alerts, setAlerts] = React.useState<AlertInfo[]>([]);
  const [ creationError, setCreationError ] = useState<boolean>(false);

  const onSubmit = async () => {
    // Create a team
    setCreationLoading(true)
    const data = {
      "team_name": newTeamName,
      "description": newTeamDesc
    }

    createTeam(data).then(response => {
      if(response.code == 200) {
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
      console.log(error)
    })

    
    // Create a repo
    try{
      const data = {
        git_organization: newOrgName,
        repository_name : newRepoName,
        jobs: {
          github_actions: {
            monitor: false
          }
        },
        artifacts: [],
        team_name: newTeamName
      }
      let response = await createRepository(data)
      if(response.code == 200) {
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
    setCreationLoading(false)
    if(!creationError) {
      setAlerts(prevAlertInfo => [...prevAlertInfo, {
        title: 'You resources have created. You will be redirected to home in a few seconds.',
        variant: AlertVariant.info,
        key: "team-not-created"
      }]);
      setTimeout(()=>{ history.push("/home/overview"); window.location.reload(); }, 3000)
    }
  };

  const onNext = (id) => {
    setState(id.id);
  };

  const onBack = (id) => {
    setState(id.id);
  };

  const onClear = () => {
    if(stepIdReached ==  "team"){
      setNewTeamName("")
      setNewTeamDesc("")
    }
    if(stepIdReached ==  "repo"){
      setNewOrgName("")
      setNewRepoName("")
    }
  };

  const TeamData = (
      <div className={'pf-u-m-lg'} >
        <Form>
          <FormGroup label="Team Name" isRequired fieldId="team-name" helperText="Include the name for your team">
            <TextInput value={newTeamName} type="text" onChange={(value)=>{setNewTeamName(value)}} aria-label="text input example" />
          </FormGroup>
          <FormGroup label="Description" fieldId='team-description' helperText="Include a description for your team">
            <TextArea value={newTeamDesc} onChange={(value)=>{setNewTeamDesc(value)}} aria-label="text area example" />
          </FormGroup>
        </Form>
      </div>
    )

  const AddRepo = (
      <div className={'pf-u-m-lg'} >
        <Form>
          <FormGroup label="Repository Name" isRequired fieldId="repo-name" helperText="Add a repository">
            <TextInput value={newRepoName} type="text" onChange={value => setNewRepoName(value)} aria-label="text input example" />
          </FormGroup>
          <FormGroup label="Organization Name" isRequired fieldId="org-name" helperText="Specify the organization">
            <TextInput value={newOrgName} type="text" onChange={value => setNewOrgName(value)} aria-label="text input example" />
          </FormGroup>
        </Form>
      </div>
    )

  const DataReview = (
      <div>
        <Title headingLevel="h6" size="xl">Review your data</Title>
        <div style={{marginTop: '2em'}}>
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
        <div style={{marginTop: "2em"}}>
          { creationLoading && <Spinner isSVG aria-label="Contents of the basic example" /> }
          <AlertGroup isLiveRegion aria-live="polite" aria-relevant="additions text" aria-atomic="false">
            {alerts.map(({ title, variant, key }) => (
              <Alert variant={variant} isInline isPlain title={title} key={key} />
            ))}
          </AlertGroup>
        </div>
      </div>
    )

  const ValidateTeamName = () => { return newTeamName != "" }
  const ValidateRepoAndOrg = () => { return newRepoName != "" && newOrgName != "" }

  const steps = [
    { id: 'team', name: 'Team Name', component: TeamData, enableNext: ValidateTeamName() },
    { id: 'repo', name: 'Add a repository', component: AddRepo, canJumpTo: ValidateTeamName(), enableNext: ValidateRepoAndOrg()},
    {
      id: 'review',
      name: 'Review',
      component: DataReview,
      nextButtonText: 'Create',
      canJumpTo: (ValidateTeamName() && ValidateRepoAndOrg()),
      hideCancelButton: true
    }
  ];
  const title = 'Incrementally enabled wizard';

  return (
    <React.Fragment>
      <PageSection style={{backgroundColor: 'white'}} variant={PageSectionVariants.light}>
        <Wizard
          navAriaLabel={`${title} steps`}
          mainAriaLabel={`${title} content`}
          steps={steps}
          onNext={onNext}
          onBack={onBack}
          onClose={onClear}
          onSave={onSubmit}
          cancelButtonText="Clear"
          height={600}
        />
      </PageSection>
    </React.Fragment>
  );
}