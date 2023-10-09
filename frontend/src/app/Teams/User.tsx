/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
import { Dropdown, DropdownGroup, DropdownItem, DropdownSeparator, DropdownToggle, ToolbarGroup, ToolbarItem } from "@patternfly/react-core"
import { CaretDownIcon, DisconnectedIcon, MailchimpIcon, UserCircleIcon } from "@patternfly/react-icons";
import React, { useEffect, useState } from "react"
import { ReactReduxContext } from "react-redux";
import { useHistory } from "react-router-dom";

export const UserToolbarGroup = () => {
    const history = useHistory();
    const [isDropdownOpen, setDropdownOpen] = useState(false);
    const { store } = React.useContext(ReactReduxContext);
    const [username, setUsername] = React.useState<string>("");
    const [userEmail, setUserEmail] = React.useState<string>("");
    const state = store.getState();

    const onDropdownToggle = (isDropdownOpen: boolean) => {
      setDropdownOpen(isDropdownOpen);
    }

    useEffect(() => {
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

    const UserDropDownItems = [
        <DropdownGroup key="group-1">
        <DropdownItem key="group-1-plaintext" component="div" isPlainText>
          {userEmail}
        </DropdownItem>
        <DropdownSeparator key="dropdown-separator" />
        </DropdownGroup>,
        <DropdownGroup key="group-2">
          <DropdownItem onClick={LogOut} key="group-2-logout" icon={<DisconnectedIcon size="lg" color="#C9190B"></DisconnectedIcon>}>Logout</DropdownItem>
        </DropdownGroup>
    ];

    return (
        <ToolbarGroup id="toolbar-user" alignment={{ default: "alignRight" }}>
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