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

  const dropdownItems = [
    <DropdownItem key="Select Team" data-value="Select Team">Select Team</DropdownItem>,
    <DropdownItem key="Link 1" data-value="link_1">Link 1</DropdownItem>,
    <DropdownItem key="Link 2" data-value="link_2">Link 2 </DropdownItem>,
    <DropdownItem key="Link 3" data-value="link_3">Link 3</DropdownItem>,
  ];

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