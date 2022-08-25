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
} from '@patternfly/react-core';
import CaretDownIcon from '@patternfly/react-icons/dist/js/icons/caret-down-icon';
import { Context } from '@app/store/store';

export interface ITeam {
  id: string
  team_name: string
}

export const BasicMasthead = () => {

  const { state, dispatch } = useContext(Context)
  const [isDropdownOpen, setDropdownOpen] = useState(false)

  const onDropdownToggle = (isDropdownOpen: boolean) => {
    setDropdownOpen(isDropdownOpen)
  };

  const onDropdownSelect = (event: any) => {
    setDropdownOpen(!isDropdownOpen)
    dispatch({ type: "SET_TEAM", data:  event.target.dataset.value });
  };

  const dropdownItems = state.TeamsAvailable.map((team) => <DropdownItem key={team.id} data-value={team.team_name}>{team.team_name}</DropdownItem> )

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
                        {state.Team}
                      </DropdownToggle>
                    }
                    isOpen={isDropdownOpen}
                    dropdownItems={dropdownItems}
                    isFullHeight
                  />
                </ToolbarItem>
              </ToolbarGroup>
            </ToolbarContent>
          </Toolbar>
        </MastheadContent>
      </Masthead>
    );
  }