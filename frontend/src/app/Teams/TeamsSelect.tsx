import React, { useState, useContext, useEffect } from 'react';
import {
  Masthead,
  MastheadContent,
  Toolbar,
  ToolbarContent,
  ToolbarGroup,
  ToolbarItem,
  Dropdown,
  DropdownToggle,
  DropdownItem,
  Button
} from '@patternfly/react-core';
import CaretDownIcon from '@patternfly/react-icons/dist/js/icons/caret-down-icon';
import { useHistory } from 'react-router-dom';
import { teamIsNotEmpty } from '@app/utils/utils';
import ArrowRightIcon from '@patternfly/react-icons/dist/esm/icons/arrow-right-icon';
import { SignOutAltIcon } from '@patternfly/react-icons';
import { ReactReduxContext, useSelector } from 'react-redux';

export interface ITeam {
  id: string
  team_name: string
  description: string
}

export const BasicMasthead = () => {
  const history = useHistory();

  const { store } = useContext(ReactReduxContext);
  const state = store.getState();
  const dispatch = store.dispatch;
  const [isDropdownOpen, setDropdownOpen] = useState(false);
  const [dropdownItems, setDropdownItems] = useState([]);

  const onDropdownToggle = (isDropdownOpen: boolean) => {
    setDropdownOpen(isDropdownOpen);
  }

  const onDropdownSelect = (event: any) => {
    setDropdownOpen(!isDropdownOpen);
    dispatch({ type: "SET_TEAM", data: event.target.dataset.value });
  }
  function Log_out() {
    history.push('/oauth/sign_out');
    window.location.reload();
  }

  const currentTeamsAvailable = useSelector((state: any) => state.teams.TeamsAvailable);

  useEffect(() => {
    let ddi = currentTeamsAvailable.map((team) => <DropdownItem key={team.id} data-value={team.team_name}>{team.team_name}</DropdownItem>)
    setDropdownItems(ddi)
  }, [currentTeamsAvailable]);

  return (
    <Masthead id="basic-demo">
      <MastheadContent>
        <Toolbar id="toolbar" isFullHeight isStatic>
          <ToolbarContent>
            <ToolbarGroup alignment={{ default: 'alignLeft' }}>
              <ToolbarItem visibility={{ default: 'hidden', lg: 'visible' }}>
                <Dropdown
                  onSelect={onDropdownSelect}
                  toggle={
                    <DropdownToggle id="toggle-id" onToggle={onDropdownToggle} toggleIndicator={CaretDownIcon}>
                      {state.teams.Team}
                    </DropdownToggle>
                  }
                  isOpen={isDropdownOpen}
                  dropdownItems={dropdownItems}
                  isFullHeight
                />
              </ToolbarItem>
              <ToolbarItem visibility={{ default: 'hidden', lg: 'visible' }}>
                {!teamIsNotEmpty(state.teams.Team) && <Button style={{ verticalAlign: "middle", backgroundColor: "var(--pf-c-button--m-primary--Color)", color: "var(--pf-c-button--m-primary--BackgroundColor)" }} onClick={() => { history.push("/home/teams") }} type="button" width={300} variant="primary">Create your first Team <ArrowRightIcon /></Button>}
              </ToolbarItem>
            </ToolbarGroup>
          </ToolbarContent>
        </Toolbar>
        <Button onClick={Log_out} variant="link" icon={<SignOutAltIcon />} iconPosition="right">
          Log out
        </Button>
      </MastheadContent>
    </Masthead>
  );
}