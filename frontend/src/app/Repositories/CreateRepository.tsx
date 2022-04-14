import { createRepository } from "@app/utils/APIService";
import { Button, ButtonVariant, Checkbox, Form, FormGroup, Modal, ModalVariant, Popover, TextArea, TextInput } from "@patternfly/react-core";
import { HelpIcon } from "@patternfly/react-icons";
import { Context } from '@app/store/store';
import React, { useContext } from "react";

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
export const FormModal = ()=> {
    const [isModalOpen, setModalOpen] = React.useState(false);
    const [gitRepositoryValue, setGitRepositoryValue] = React.useState("");
    const [gitOrganizationValue, setGitOrganizationValue] = React.useState("");
    const [monitorGithubActions, setMonitorGithubActions] = React.useState(false);
    const [checked, setChecked] = React.useState('');
  
    const handleModalToggle = () => {
      setModalOpen(!isModalOpen);
    };
  
    const handleGitRepositoryInput = value => {
      setGitRepositoryValue(value);
    };
  
    const handleGitOrganizationInput = value => {
      setGitOrganizationValue(value);
    };
  
    const handleGithubActionsMonitor = value => {
      setMonitorGithubActions(value);
      setChecked(value);
    };

    const  onSubmit = async() => {
      try{
      const data = {
        git_organization: gitOrganizationValue,
        repository_name : gitRepositoryValue,
        jobs: {
          github_actions: {
            monitor: monitorGithubActions
          }
        },
        artifacts: []
      }
      handleModalToggle()
      await createRepository(data)
      window.location.reload();
    }
    catch (error) {
      console.log(error)
    }
  }
  
    return (
      <React.Fragment>
        <Button variant={ButtonVariant.danger} onClick={handleModalToggle}>
          Add Git Repository
        </Button>
        <Modal
         variant={ModalVariant.medium}
          title="Add new git repository"
          description="Enter a new git repository to obtain information in the quality dashboard."
          isOpen={isModalOpen}
          onClose={handleModalToggle}
          actions={[
            <Button key="create" variant="primary" form="modal-with-form-form" onClick={onSubmit}>
              Add
            </Button>,
            <Button key="cancel" variant="link" onClick={handleModalToggle}>
              Cancel
            </Button>
          ]}
        >
         <Form isHorizontal id="modal-with-form-form">
          <FormGroup 
            label="Git Organization"
            labelIcon={
              <Popover
                headerContent={
                  <div>
                  </div>
                }
                bodyContent={
                  <div>
                    Git organization name
                  </div>
                }
              >
                <button
                  type="button"
                  aria-label="More info for name field"
                  onClick={e => e.preventDefault()}
                  aria-describedby="modal-with-form-form-name"
                  className="pf-c-form__group-label-help"
                >
                  <HelpIcon noVerticalAlign />
                </button>
              </Popover>
            }
            isRequired
            fieldId="modal-with-form-form-name">
              <TextInput
                isRequired
                type="email"
                id="modal-with-form-form-name"
                name="modal-with-form-form-name"
                value={gitOrganizationValue}
                onChange={handleGitOrganizationInput}
              />
          </FormGroup>
          <FormGroup 
            label="Repository name"
            labelIcon={
              <Popover
                headerContent={
                  <div>
                    The repository name
                  </div>
                }
                bodyContent={
                  <div>
                    An valid github repository bane
                  </div>
                }
              >
                <button
                  type="button"
                  aria-label="More info for e-mail field"
                  onClick={e => e.preventDefault()}
                  aria-describedby="modal-with-form-form-email"
                  className="pf-c-form__group-label-help"
                >
                  <HelpIcon noVerticalAlign />
                </button>
              </Popover>
            }
            isRequired
            fieldId="modal-with-form-form-email">
              <TextInput
                isRequired
                type="email"
                id="modal-with-form-form-email"
                name="modal-with-form-form-email"
                value={gitRepositoryValue}
                onChange={handleGitRepositoryInput}
              />
          </FormGroup>
          <FormGroup label="Monitor CI Jobs" isRequired isStack hasNoPaddingTop fieldId={''}>
            <Checkbox label="Github Actions" id="alt-form-checkbox-1" name="alt-form-checkbox-1" value={monitorGithubActions} onChange={handleGithubActionsMonitor} isChecked={checked}/>
          </FormGroup>
          <FormGroup label="Code Coverage" isRequired isStack hasNoPaddingTop fieldId={''}>
            <Checkbox label="codecov.io" id="alt-form-checkbox-2" name="alt-form-checkbox-2" />
          </FormGroup>
          <FormGroup 
            label="Quay.io Artifacts"
            labelIcon={
              <Popover
                headerContent={
                  <div>
                    Quay.io artifacts related to the repository
                  </div>
                }
                bodyContent={
                  <div>
                    If the repository contain more than one artifact add with coma separated:
                    quay.io/flacatus:repo1,quay.io/flacatus:repo2
                  </div>
                }
              >
                <button
                  type="button"
                  aria-label="More info for address field"
                  onClick={e => e.preventDefault()}
                  aria-describedby="modal-with-form-form-address"
                  className="pf-c-form__group-label-help"
                >
                  <HelpIcon noVerticalAlign />
                </button>
              </Popover>
            }
            isRequired
            fieldId="modal-with-form-form-address">
            <TextArea
              name="horizontal-form-exp"
              id="horizontal-form-exp"
            />
          </FormGroup>
         </Form>
        </Modal>
      </React.Fragment>
    );
  };

