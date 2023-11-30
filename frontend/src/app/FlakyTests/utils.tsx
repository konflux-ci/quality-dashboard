import React from 'react';
import { Dropdown, DropdownToggle, DatePicker, Spinner } from '@patternfly/react-core';

export const SpinnerBasic: React.FunctionComponent<{ isLoading: boolean }> = ({ isLoading }) => {
    return (
        <React.Fragment>
            {isLoading &&
                <div className='spinner-loading'>
                    <Spinner isSVG aria-label="Contents of the basic example" />
                </div>
            }
        </React.Fragment>
    )
};

export const DropdownBasic: React.FunctionComponent<{ toggles, onSelect, selected, placeholder }> = ({ toggles, onSelect, selected, placeholder }) => {
    const [isOpen, setIsOpen] = React.useState(false);

    const onToggle = (isOpen: boolean) => {
        setIsOpen(isOpen);
    };

    const onFocus = () => {
        const element = document.getElementById('toggle-basic');
        element?.focus();
    };

    const onItemSelect = (e) => {
        setIsOpen(false);
        onFocus();
        onSelect(e.target.name)
    };

    return (
        <Dropdown
            onSelect={onItemSelect}
            toggle={<DropdownToggle onToggle={onToggle}>{selected == '' ? placeholder : selected}</DropdownToggle>}
            isOpen={isOpen}
            dropdownItems={toggles}
        />
    );
};

export const DatePickerMinMax: React.FunctionComponent<{ selectedDate: string | undefined, onChange: (value: string, date: Date, name: string) => void, name: string }> = ({ selectedDate, onChange, name }) => {
    const onDateChange = (e, value, date) => {
        onChange(value, date, name)
    }
    return <DatePicker name={name} onChange={onDateChange} value={selectedDate?.split('T')[0]} />;
};