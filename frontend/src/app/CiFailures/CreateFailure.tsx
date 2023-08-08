import { bugExists, createFailure } from '@app/utils/APIService';
import {
    Button,
    Form,
    FormGroup,
    Modal,
    ModalVariant,
    Popover,
    TextInput,
    Alert,
} from '@patternfly/react-core';
import { HelpIcon } from '@patternfly/react-icons';
import React, { useContext, SetStateAction, useEffect } from 'react';
import { teamIsNotEmpty } from '@app/utils/utils';
import { useHistory } from 'react-router-dom';
import { ReactReduxContext } from 'react-redux';
import { formatDate, getRangeDates } from '@app/Reports/utils';

interface IModalContext {
    isModalOpen: IModalContextMember;
    isEdit: IModalContextMember;
    handleModalToggle;
    data: IModalContextMember;
}

interface IModalContextMember {
    value: any;
    set: SetStateAction<any>;
}
export const useDefaultModalContextState = () => {
    const [isModalOpen, setModalOpen] = React.useState(false);
    const [isEdit, setEditRepo] = React.useState(false);
    const [data, setData] = React.useState({});
    const defaultModalContext = useModalContext();

    defaultModalContext.isModalOpen = { set: setModalOpen, value: isModalOpen };
    defaultModalContext.isEdit = { set: setEditRepo, value: isEdit };
    defaultModalContext.data = { set: setData, value: data };
    defaultModalContext.handleModalToggle = (edit: Boolean, data: any) => {
        defaultModalContext.isModalOpen.set(!defaultModalContext.isModalOpen.value);
        if (edit == true) {
            defaultModalContext.isEdit.set(true);
            defaultModalContext.data.set(data);
        } else {
            defaultModalContext.isEdit.set(false);
            defaultModalContext.data.set({ jira_key: '', error_message: '' });
        }
    };
    return defaultModalContext;
};

export const ModalContext = React.createContext<IModalContext>({
    isModalOpen: { set: undefined, value: false },
    isEdit: { set: undefined, value: false },
    handleModalToggle: () => { },
    data: { set: undefined, value: false },
});

export const useModalContext = () => {
    return useContext(ModalContext);
};

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
export const FormModal = () => {
    const { store } = React.useContext(ReactReduxContext);
    const state = store.getState();
    const modalContext = useModalContext();
    const history = useHistory();
    const [jiraKeyValue, setJiraKeyValue] = React.useState('');
    const [errorMessageValue, setErrorMessageValue] = React.useState('');
    const params = new URLSearchParams(window.location.search);
    type validate = 'success' | 'warning' | 'error' | 'default';
    const [jiraKeyValidated, setJiraKeyValidated] = React.useState<validate>('error');
    const [msgValidated, setMsgValidated] = React.useState<validate>('error');

    const handleJiraKeyInput = async (value) => {
        setJiraKeyValidated('error');
        setJiraKeyValue(value);
        if (value != "" && value != undefined) {
            bugExists(value, state.teams.Team).then(res => {
                if (res) {
                    setJiraKeyValidated('success');
                }
            })
        } else {
            setJiraKeyValidated('error');
        }
    };

    const handleErrorMessageInput = (value) => {
        setErrorMessageValue(value);
        if (value != "") {
            setMsgValidated('success');
        } else {
            setMsgValidated('error');
        }
    };

    interface LoadingPropsType {
        spinnerAriaValueText: string;
        spinnerAriaLabelledBy?: string;
        spinnerAriaLabel?: string;
        isLoading: boolean;
    }

    const [isPrimaryLoading, setIsPrimaryLoading] = React.useState<boolean>(false);
    const primaryLoadingProps = {} as LoadingPropsType;
    primaryLoadingProps.spinnerAriaValueText = 'Loading';
    primaryLoadingProps.spinnerAriaLabelledBy = 'primary-loading-button';
    primaryLoadingProps.isLoading = isPrimaryLoading;


    const onCreateSubmit = async () => {
        setIsPrimaryLoading(!isPrimaryLoading);

        try {
            await createFailure(state.teams.Team, jiraKeyValue, errorMessageValue);
            setIsPrimaryLoading(!isPrimaryLoading);
            modalContext.handleModalToggle();

            const rangeDateTime = getRangeDates(15);
            const start_date = formatDate(rangeDateTime[0]);
            const end_date = formatDate(rangeDateTime[1]);

            history.push(
                '/home/rhtapbugs-impact?team=' +
                params.get('team') +
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

    //   if (modalContext.isEdit.value) {
    //         setJiraKeyValue(modalContext.data.value.jira_key);
    //         // setErrorMessageValue(modalContext.data.value.error_message);
    //     }

    useEffect(() => {
        setJiraKeyValue(modalContext.data.value.jira_key);
        handleJiraKeyInput(modalContext.data.value.jira_key);
        setErrorMessageValue(modalContext.data.value.error_message);
        handleErrorMessageInput(modalContext.data.value.error_message);
    }, [modalContext.data.value.jira_key, modalContext.data.value.error_message]);

    return (
        <React.Fragment>
            <Modal
                variant={ModalVariant.medium}
                title={!modalContext.isEdit.value ? 'Add a RHTAPBUG' : 'Update'}
                description={
                    !modalContext.isEdit.value
                        ? 'Track the impact of a RHTAPBUG'
                        : ''
                }
                isOpen={modalContext.isModalOpen.value}
                onClose={modalContext.handleModalToggle}
                actions={[
                    <Button
                        key="create"
                        variant="primary"
                        form="modal-with-form-form"
                        onClick={onCreateSubmit}
                        isDisabled={msgValidated == "error" || jiraKeyValidated == "error"}
                        {...primaryLoadingProps}
                    >
                        {!modalContext.isEdit.value ? 'Add' : 'Update'}
                    </Button>,
                    <Button key="cancel" variant="link" onClick={modalContext.handleModalToggle}>
                        Cancel
                    </Button>,
                ]}
            >
                <Form isHorizontal id="modal-with-form-form">
                    <FormGroup
                        label="Jira Key"
                        labelIcon={
                            <Popover headerContent={<div></div>} bodyContent={<div>Add a valid Jira Key that refers to a ci-fail issue</div>}>
                                <button
                                    type="button"
                                    aria-label="More info for name field"
                                    onClick={(e) => e.preventDefault()}
                                    aria-describedby="modal-with-form-form-name"
                                    className="pf-c-form__group-label-help"
                                >
                                    <HelpIcon noVerticalAlign />
                                </button>
                            </Popover>
                        }
                        isRequired
                        fieldId="modal-with-form-form-name"
                        validated={jiraKeyValidated}
                        helperTextInvalid="Must be a valid JIRA key"
                    >
                        <TextInput
                            validated={jiraKeyValidated}
                            isReadOnly={modalContext.isEdit.value}
                            isRequired
                            type="email"
                            id="modal-with-form-form-name"
                            name="modal-with-form-form-name"
                            value={jiraKeyValue}
                            onChange={handleJiraKeyInput}
                        />
                    </FormGroup>
                    <FormGroup
                        label="Error message"
                        labelIcon={
                            <Popover
                                headerContent={<div>Error message</div>}
                                bodyContent={<div>Add a error message</div>}
                            >
                                <button
                                    type="button"
                                    aria-label="More info for e-mail field"
                                    onClick={(e) => e.preventDefault()}
                                    aria-describedby="modal-with-form-form-email"
                                    className="pf-c-form__group-label-help"
                                >
                                    <HelpIcon noVerticalAlign />
                                </button>
                            </Popover>
                        }
                        isRequired
                        fieldId="modal-with-form-form-email"
                        validated={msgValidated}
                        helperTextInvalid="Must be a error message"
                    >
                        <TextInput
                            validated={msgValidated}
                            isRequired
                            type="email"
                            id="modal-with-form-form-email"
                            name="modal-with-form-form-email"
                            value={errorMessageValue}
                            onChange={value => handleErrorMessageInput(value)}
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
                </Form>
            </Modal>
        </React.Fragment>
    );
};

