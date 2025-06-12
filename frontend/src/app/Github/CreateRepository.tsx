import { checkGithubRepositoryExists, checkGithubRepositoryUrl, createRepository } from '@app/utils/APIService';
import {
  Button,
  Checkbox,
  Form,
  FormGroup,
  Modal,
  ModalVariant,
  Popover,
  TextArea,
  TextInput,
  Alert,
} from '@patternfly/react-core';
import { HelpIcon } from '@patternfly/react-icons';
import React, { useContext, SetStateAction, useEffect, useState } from 'react';
import { teamIsNotEmpty } from '@app/utils/utils';
import { useHistory } from 'react-router-dom';
import { ReactReduxContext } from 'react-redux';
import { formatDateTime, getRangeDates } from '@app/Reports/utils';
import { ghRegex, githubRegExp } from '@app/Teams/TeamsOnboarding';

interface IModalContext {
  isModalOpen: IModalContextMember;
  isEditRepo: IModalContextMember;
  handleModalToggle;
  data: IModalContextMember;
}

interface IModalContextMember {
  value: any;
  set: SetStateAction<any>;
}

export interface LoadingPropsType {
  spinnerAriaValueText: string;
  spinnerAriaLabelledBy?: string;
  spinnerAriaLabel?: string;
  isLoading: boolean;
}

export const ModalContext = React.createContext<IModalContext>({
  isModalOpen: { set: undefined, value: false },
  isEditRepo: { set: undefined, value: false },
  handleModalToggle: () => { },
  data: { set: undefined, value: false },
});

export const useModalContext = () => {
  return useContext(ModalContext);
};

export const useDefaultModalContextState = () => {
  const [isModalOpen, setModalOpen] = React.useState(false);
  const [isEditRepo, setEditRepo] = React.useState(false);
  const [data, setData] = React.useState({});
  const defaultModalContext = useModalContext();

  defaultModalContext.isModalOpen = { set: setModalOpen, value: isModalOpen };
  defaultModalContext.isEditRepo = { set: setEditRepo, value: isEditRepo };
  defaultModalContext.data = { set: setData, value: data };
  defaultModalContext.handleModalToggle = (edit: Boolean, data: any) => {
    defaultModalContext.isModalOpen.set(!defaultModalContext.isModalOpen.value);
    if (edit == true) {
      defaultModalContext.isEditRepo.set(true);
      defaultModalContext.data.set(data);
    } else {
      defaultModalContext.isEditRepo.set(false);
    }
  };
  return defaultModalContext;
};

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
export const FormModal: React.FunctionComponent = () => {
  const modalContext = useModalContext();
  const history = useHistory();
  const [gitRepositoryValue, setGitRepositoryValue] = React.useState('');
  const [artifactsValue, setArtifactsValue] = React.useState<string[]>([]);
  const [gitOrganizationValue, setGitOrganizationValue] = React.useState('');
  const [monitorGithubActions, setMonitorGithubActions] = React.useState(false);
  const [checked, setChecked] = React.useState('');
  const params = new URLSearchParams(window.location.search);
  const { store } = useContext(ReactReduxContext);
  const state = store.getState();
  const dispatch = store.dispatch;
  type validate = 'success' | 'warning' | 'error' | 'default';
  const [githubUrlValidated, setGithubUrlValidated] = React.useState<validate>('error');
  const [githubUrl, setGithubUrl] = useState<string>("");
  const [helperText, setHelperText] = useState<string>("Enter a GitHub Repository URL");

  const handleGitRepositoryInput = (value) => {
    setGitRepositoryValue(value);
  };

  const handleGitOrganizationInput = (value) => {
    setGitOrganizationValue(value);
  };

  const handleArtifactsValue = (value: string) => {
    const artifactsArray = value ? value.split(',').map(item => item.trim()) : [];
    setArtifactsValue(artifactsArray);
  };

  const handleGithubActionsMonitor = (value) => {
    setMonitorGithubActions(value);
    setChecked(value);
  };

  const onUpdateSubmit = async () => {
    modalContext.handleModalToggle();
    window.location.reload();
  };

  const [isPrimaryLoading, setIsPrimaryLoading] = React.useState<boolean>(false);
  const primaryLoadingProps = {} as LoadingPropsType;
  primaryLoadingProps.spinnerAriaValueText = 'Loading';
  primaryLoadingProps.spinnerAriaLabelledBy = 'primary-loading-button';
  primaryLoadingProps.isLoading = isPrimaryLoading;

  const onCreateSubmit = async () => {
    setIsPrimaryLoading(!isPrimaryLoading);
    try {
      const data = {
        git_organization: gitOrganizationValue,
        repository_name: gitRepositoryValue,
        jobs: {
          github_actions: {
            monitor: monitorGithubActions,
          },
        },
        artifacts: artifactsValue,
        team_name: state.teams.Team,
      };

      await createRepository(data);
      setIsPrimaryLoading(!isPrimaryLoading);
      modalContext.handleModalToggle();

      const rangeDateTime = getRangeDates(90);
      const start_date = formatDateTime(rangeDateTime[0]);
      const end_date = formatDateTime(rangeDateTime[1]);

      history.push(
        '/home/github?team=' +
        params.get('team') +
        '&organization=' +
        gitOrganizationValue +
        '&repository=' +
        gitRepositoryValue +
        '&start=' +
        start_date +
        '&end=' +
        end_date
      );
      window.location.reload();
    } catch (error) {
      console.log(error);
    }
  };

  const onSubmit = async () => {
    !modalContext.isEditRepo.value ? onCreateSubmit() : onUpdateSubmit();
  };

  useEffect(() => {
    if (modalContext.isEditRepo.value) {
      setGitRepositoryValue(modalContext.data.value.repository_name);
      setGitOrganizationValue(modalContext.data.value.git_organization);
    }
  });

const handleGithub = async (value: string) => {
   const githubRegExp = /^https:\/\/github\.com\/([^/]+)\/([^/]+)\/?$/;

    setHelperText('')
    setGithubUrl(value);
    setGithubUrlValidated('error');

    // Use .match() to capture the owner and repo from the URL
    // The regex should have capturing groups: e.g., /^https:\/\/github\.com\/([^/]+)\/([^/]+)/
    const match = value.match(githubRegExp);

    // check that the URL format is valid and we have our captured values
    if (match) {
      const owner = match[1];      // The first captured group (e.g., "org")
      const repository = match[2]; // The second captured group (e.g., "repo")

      // check that gh repo was not already added
      const resp = await checkGithubRepositoryExists(owner, repository)
      if (resp != undefined && resp.code == 200) {
        const team = resp.data as string
        setHelperText('Already exists in ' + team + ' team')
      } else {
        await checkGithubRepositoryUrl(owner, repository).then((data: any) => {
          if (data != undefined && data.code == 200) {
            setGithubUrlValidated('success');
            // save repo and org
            setGitOrganizationValue(owner);
            setGitRepositoryValue(repository);
            setHelperText('')
          } else {
            setHelperText('Something went wrong. Probably URL is incorrect.')
          }
        })
      }
    } else {
      setHelperText('Must match the regex `' + ghRegex +'`')
    }
  };


  const [age, setAge] = React.useState('Five');
  const [validated, setValidated] = React.useState<validate>('error');

  const handleAgeChange = (age: string) => {
    setAge(age);
    if (age === '') {
      setValidated('default');
    } else if (/^\d+$/.test(age)) {
      setValidated('success');
    } else {
      setValidated('error');
    }
  };

  return (
    <Modal
      variant={ModalVariant.medium}
      title={!modalContext.isEditRepo.value ? 'Add new GitHub Repository' : 'Update GitHub Repository'}
      description={
        !modalContext.isEditRepo.value
          ? 'Enter a new GitHub Repository to obtain information in the quality studio.'
          : ''
      }
      isOpen={modalContext.isModalOpen.value}
      onClose={modalContext.handleModalToggle}
      actions={[
        <Button
          key="create"
          variant="primary"
          form="modal-with-form-form"
          onClick={onSubmit}
          isDisabled={githubUrlValidated == "error"}
          {...primaryLoadingProps}
        >
          {!modalContext.isEditRepo.value ? 'Add' : 'Update'}
        </Button>,
        <Button key="cancel" variant="link" onClick={modalContext.handleModalToggle}>
          Cancel
        </Button>,
      ]}
    >
      <Form isHorizontal id="modal-with-form-form">
        <FormGroup isRequired label="GitHub Repository URL" fieldId="repo-name" helperText={helperText}>
          <TextInput
            isRequired
            validated={githubUrlValidated}
            value={githubUrl}
            type="text"
            onChange={handleGithub}
            aria-label="text input example"
            placeholder="Add a GitHub repository"
          />
        </FormGroup>
        <FormGroup label="Team" isRequired isStack hasNoPaddingTop fieldId={''}>
          {teamIsNotEmpty(state.teams.Team) ? (
            <TextInput
              isReadOnly={true}
              isRequired
              type="text"
              id="modal-with-form-form-team"
              name="modal-with-form-form-team"
              value={state.teams.Team}
            />
          ) : (
            <div>
              <Button
                onClick={() => {
                  history.push('/home/teams');
                }}
                type="button"
                width={300}
              >
                Create your first Team
              </Button>
              <Alert
                style={{ marginTop: '1em' }}
                variant="danger"
                isInline
                isPlain
                title="You need to create a team before adding a repository"
              />
            </div>
          )}
        </FormGroup>
        <FormGroup label="Monitor CI Jobs" isRequired isStack hasNoPaddingTop fieldId={''}>
          <Checkbox
            label="GitHub Actions"
            id="alt-form-checkbox-1"
            name="alt-form-checkbox-1"
            value={String(monitorGithubActions)}
            onChange={handleGithubActionsMonitor}
            isChecked={Boolean(checked)}
          />
        </FormGroup>
        <FormGroup label="Code Coverage" isRequired isStack hasNoPaddingTop fieldId={''}>
          <Checkbox label="codecov.io" id="alt-form-checkbox-2" name="alt-form-checkbox-2" />
        </FormGroup>
        <FormGroup
          label="Quay.io Artifacts"
          labelIcon={
            <Popover
              headerContent={<div>Quay.io artifacts related to the repository</div>}
              bodyContent={
                <div>
                  If the repository contain more than one artifact add with coma separated:
                  quay.io/team-storage/repo1,quay.io/team-storage/repo2
                </div>
              }
            >
              <button
                type="button"
                aria-label="More info for address field"
                onClick={(e) => e.preventDefault()}
                aria-describedby="modal-with-form-form-address"
                className="pf-c-form__group-label-help"
              >
                <HelpIcon noVerticalAlign />
              </button>
            </Popover>
          }
          isRequired
          fieldId="modal-with-form-form-address"
        >
            <TextArea
            name="horizontal-form-exp"
            id="horizontal-form-exp"
            value={artifactsValue.join(', ')} // Set the component's value from state
            onChange={handleArtifactsValue}   // Call the handler when the user types
          />
        </FormGroup>
      </Form>
    </Modal>
  );
};
