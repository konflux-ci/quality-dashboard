/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
import { createUser, listUsers } from "@app/utils/APIService";
import { Button, Dropdown, DropdownGroup, DropdownItem, DropdownSeparator, DropdownToggle, Form, FormGroup, FormSection, FormSelect, FormSelectOption, Modal, ModalVariant, Popover, TextInput, ToolbarGroup, ToolbarItem } from "@patternfly/react-core"
import { CaretDownIcon, CogIcon, DisconnectedIcon, UserCircleIcon } from "@patternfly/react-icons";
import React, { useEffect, useState } from "react"
import { ReactReduxContext, useSelector } from "react-redux";
import { useHistory } from "react-router-dom";

export interface TeamsConfiguration {
    default_team: string;
}

export interface UserConfig {
    teams_configuration: TeamsConfiguration;
}

const getUserConfig = (defaultTeam) => {
    const cfg: UserConfig = {
        teams_configuration: {
            default_team: defaultTeam,
        }
    };
    const config = JSON.stringify(cfg);

    return config
}

export const UserToolbarGroup = () => {
    const history = useHistory();
    const [isDropdownOpen, setDropdownOpen] = useState(false);
    const { store } = React.useContext(ReactReduxContext);
    const [username, setUsername] = React.useState<string>("");
    const [userEmail, setUserEmail] = React.useState<string>("");
    const state = store.getState();
    const [isModalOpen, setIsModalOpen] = React.useState(false);
    const [defaultTeam, setDefaultTeam] = React.useState('');
    const currentTeamsAvailable = useSelector((state: any) => state.teams.TeamsAvailable);
    const redux_dispatch = store.dispatch;

    const confirm = async () => {
        // get user config
        const config = getUserConfig(defaultTeam)
        // update user config
        await createUser(userEmail, config)
        redux_dispatch({ type: "SET_USER_CONFIG", data: config });

        setIsModalOpen(!isModalOpen);
        window.location.reload();
    };

    const handleModalToggle = () => {
        setIsModalOpen(!isModalOpen);
    };

    const handleDefaultTeamChange = (value: string) => {
        setDefaultTeam(value);
    };

    const onDropdownToggle = (isDropdownOpen: boolean) => {
        setDropdownOpen(isDropdownOpen);
    }

    listUsers()

    useEffect(() => {
        setDefaultTeam("n/a")
        if (state.auth.USER_CONFIG != "") {
            var userConfig = JSON.parse(state.auth.USER_CONFIG) as UserConfig;
            const userConfigDefaultTeam = userConfig.teams_configuration.default_team
            if (userConfigDefaultTeam != "") {
                setDefaultTeam(userConfigDefaultTeam)
            }
        }

        try {
            const userClaims = JSON.parse(window.atob(state.auth.IDT.split('.')[1]))
            setUsername(userClaims.name)
            setUserEmail(userClaims.email)
        } catch (error) {
            history.push('/login');
        }
    }, []);

    const LogOut = () => {
        localStorage.clear()
        history.push('/login');
        window.location.reload();
    }

    const openSettings = () => {
        setDropdownOpen(false);
        setIsModalOpen(!isModalOpen);
    };

    const UserDropDownItems = [
        <DropdownGroup key="group-1">
            <DropdownItem key="group-1-plaintext" component="div" isPlainText>
                {userEmail}
            </DropdownItem>
            <DropdownSeparator key="dropdown-separator" />
        </DropdownGroup>,
        <DropdownGroup key="group-2">
            <DropdownItem onClick={openSettings} key="group-2-settings" icon={<CogIcon size="lg"></CogIcon>}>Settings</DropdownItem>
            <DropdownSeparator key="dropdown-separator" />
        </DropdownGroup>,
        <DropdownGroup key="group-3">
            <DropdownItem onClick={LogOut} key="group-3-logout" icon={<DisconnectedIcon size="lg" color="#C9190B"></DisconnectedIcon>}>Logout</DropdownItem>
        </DropdownGroup>
    ];

    const UserModal = () => {
        return (
            <React.Fragment>
                <Modal
                    variant={ModalVariant.medium}
                    title="User Settings"
                    isOpen={isModalOpen}
                    onClose={handleModalToggle}
                    actions={[
                        <Button key="confirm" variant="primary" onClick={confirm}>
                            Confirm
                        </Button>,
                        <Button key="cancel" variant="link" onClick={handleModalToggle}>
                            Cancel
                        </Button>
                    ]}
                    ouiaId="BasicModal"
                >
                    <Form isHorizontal id="modal-with-form-form">
                        <FormSection title="Team's configuration" titleElement="h2">
                            <FormGroup label="Default team" fieldId="horizontal-form-title">
                                <FormSelect
                                    value={defaultTeam}
                                    onChange={handleDefaultTeamChange}
                                    id="horizontal-form-title"
                                    name="horizontal-form-title"
                                    aria-label="Your title"
                                >
                                    {currentTeamsAvailable.map((option, index) => (
                                        <FormSelectOption key={index} value={option.team_name} label={option.team_name} />
                                    ))}
                                    <FormSelectOption key={currentTeamsAvailable.length + 1} value={"n/a"} label={"n/a"} />
                                </FormSelect>
                            </FormGroup>
                        </FormSection>
                    </Form>
                </Modal>
            </React.Fragment>)
    }

    return (
        <ToolbarGroup id="toolbar-user" alignment={{ default: "alignRight" }}>
            <UserModal></UserModal>
            <ToolbarItem visibility={{ default: 'hidden', lg: 'visible' }}>
                <Dropdown
                    toggle={
                        <DropdownToggle
                            id="toggle-id"
                            onToggle={onDropdownToggle}
                            toggleIndicator={CaretDownIcon}
                            icon={<UserCircleIcon size="lg"></UserCircleIcon>}
                        >
                            {username}
                        </DropdownToggle>
                    }
                    isOpen={isDropdownOpen}
                    dropdownItems={UserDropDownItems}
                    isFullHeight
                />
            </ToolbarItem>
        </ToolbarGroup>
    )
}