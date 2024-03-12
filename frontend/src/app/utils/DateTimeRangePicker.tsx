import { Flex, FlexItem, InputGroup, DatePicker, isValidDate, TimePicker, Select, SelectOption, Button, Popover, InputGroupText } from '@patternfly/react-core';
import { SearchIcon, CalendarAltIcon } from '@patternfly/react-icons';
import React, { useEffect, useState } from 'react';
import { customizedFormatDate, customizedFormatDateTime, getDateFormat, getRangeDateTime, getRangeDates, parseCustomizedDate, ranges } from '../Reports/utils';
import moment from 'moment';

export const DateTimeRangePicker = (props) => {
    const [from, setFrom] = React.useState(props.startDate);
    const [to, setTo] = React.useState(props.endDate);
    const [quickRangeToggle, setQuickRangeToggle] = useState(false);
    const [quickRange, setQuickRange] = useState("");
    const popoverRef = React.useRef<HTMLButtonElement>(null);
    const [isVisible, setIsVisible] = React.useState(false);

    useEffect(
        () => {
            setFrom(props.startDate);
            setTo(props.endDate);
        },
        [props.startDate, props.endDate],
    );

    const setDateTimeRange = (range: Date[]) => {
        setFrom(range[0])
        setTo(range[1])
    }

    const clearQuickRange = () => {
        setQuickRange("");
        setQuickRangeToggle(false);
    }

    const setQuickRangeOnChange = (event, selection, isPlaceholder) => {
        if (isPlaceholder) {
            clearQuickRange()
        }
        else {
            if (ranges[selection].days > 0) {
                const range = getRangeDates(ranges[selection].days)
                setDateTimeRange(range)
            }
            if (ranges[selection].hours > 0) {
                const range = getRangeDateTime(ranges[selection].hours)
                setDateTimeRange(range)
            }
            setQuickRange(ranges[selection].type);
            setQuickRangeToggle(false);
        }
    };

    const toValidator = date => {
        return date.getTime() >= from.getTime() ? '' : 'To date must be after from date';
    };
    
    const customdateValidator = date => {
        return moment(date, getDateFormat(), true).isValid() ? '' : 'Invalid date';
    }

    const onFromDateChange = (_event, inputDate, newFromDate) => {
        clearQuickRange()
        
        if (isValidDate(from) && isValidDate(newFromDate) && inputDate === customizedFormatDate(newFromDate)) {
            newFromDate.setHours(from.getHours(), to.getMinutes());
        }
        if (isValidDate(newFromDate) && inputDate === customizedFormatDate(newFromDate)) {
            setFrom(parseCustomizedDate(newFromDate));
        }
    };

    const onFromTimeChange = (_event, time, hour, minute) => {
        clearQuickRange()
        if (isValidDate(from) && time.length == 5) {
            const updatedFromDate = new Date(from);
            updatedFromDate.setHours(hour, minute);
            setFrom(updatedFromDate);
        }
    };

    const onToDateChange = (_event, inputDate, newToDate) => {
        clearQuickRange()
        if (isValidDate(to) && isValidDate(newToDate) && inputDate === customizedFormatDate(newToDate)) {
            newToDate.setHours(to.getHours(), to.getMinutes());
        }
        if (isValidDate(newToDate) && inputDate === customizedFormatDate(newToDate)) {
            setTo(newToDate);
        }
    };

    const onToTimeChange = (_event, time, hour, minute) => {
        clearQuickRange()
        if (isValidDate(to) && time.length == 5) {
            const updatedToDate = new Date(to);
            updatedToDate.setHours(hour, minute);
            setTo(updatedToDate);
        }
    };

    function wrapHandleChange(event) {
        props.handleChange(event, from, to);
        setIsVisible(false)
    }

    const open = () => {
        setIsVisible(true)
        clearQuickRange()
    }

    return (
        <div>
            <Button variant="secondary" ref={popoverRef} style={{background: 'white'}}>
                <InputGroupText id="searchFrom">
                    <CalendarAltIcon style={{ marginRight: "5px" }} />
                    {customizedFormatDateTime(from) + " to " + customizedFormatDateTime(to)}
                </InputGroupText>
            </Button>
            <Popover
                aria-label={customizedFormatDateTime(from) + " " + customizedFormatDateTime(to)}
                hasAutoWidth={true}
                isVisible={isVisible}
                shouldOpen={() => open()}
                shouldClose={() => setIsVisible(false)}
                flipBehavior={["bottom"]}
                headerContent={<div>Select date time range</div>}
                position="auto"
                bodyContent={
                    <Flex direction={{ default: 'column' }}>
                            <Flex direction={{ default: 'column' }}>
                                <FlexItem>
                                    <Select placeholderText="Search quick ranges" isOpen={quickRangeToggle} onToggle={setQuickRangeToggle} selections={quickRange} onSelect={setQuickRangeOnChange} aria-label="Select Input" toggleIcon={<SearchIcon />}>
                                        {ranges.map((value, index) => (
                                            <SelectOption key={index} value={index}>{value.type}</SelectOption>
                                        ))}
                                    </Select>
                                </FlexItem>
                                <FlexItem>
                                    From
                                    <InputGroup>
                                        <DatePicker
                                            value={isValidDate(from) ? customizedFormatDate(from) : from.toString()}
                                            onChange={onFromDateChange}
                                            aria-label="Start date"
                                            placeholder={getDateFormat()}
                                            dateParse={parseCustomizedDate}
                                            dateFormat={customizedFormatDate}
                                        />
                                        <TimePicker
                                            time={from}
                                            aria-label="Start time"
                                            style={{ width: '150px' }}
                                            is24Hour
                                            onChange={onFromTimeChange}
                                        />
                                    </InputGroup>
                                </FlexItem>
                                <FlexItem>
                                    To
                                    <InputGroup>
                                        <DatePicker
                                            value={isValidDate(to) ? customizedFormatDate(to) : to.toString()}
                                            onChange={onToDateChange}
                                            validators={[customdateValidator, toValidator]}
                                            aria-label="End date"
                                            placeholder={getDateFormat()}
                                            dateParse={parseCustomizedDate}
                                            dateFormat={customizedFormatDate}
                                        />
                                        <TimePicker
                                            time={to}
                                            aria-label="End time"
                                            style={{ width: '150px' }}
                                            is24Hour
                                            onChange={onToTimeChange}
                                            isDisabled={!isValidDate(from)}
                                        />
                                    </InputGroup>
                                </FlexItem>
                                <FlexItem>
                                    <Button variant="primary" isSmall onClick={wrapHandleChange}>
                                        Apply date time range
                                    </Button>
                                </FlexItem>
                            </Flex>
                        </Flex>
                }
                reference={popoverRef}
            />
        </div>
    );
}