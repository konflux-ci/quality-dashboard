import React from 'react';
import { CodeBranchIcon, FilterIcon, UserIcon } from '@patternfly/react-icons';
import {
    Card,
    CardTitle,
    CardBody,
    Text,
    TextVariants,
    Bullseye,
    Icon,
    FlexItem,
    Divider,
    Flex,
} from '@patternfly/react-core';


export const InfoBanner = () => (
    <Card isLarge>
        <CardBody style={{ paddingLeft: '16px' }}>
            <Flex
                justifyContent={{ default: 'justifyContentSpaceEvenly' }}
                flexWrap={{ default: 'nowrap' }}
                direction={{ default: 'column', sm: 'row' }}
            >
                <FlexItem flex={{ default: 'flex_1' }}>
                    <Card isPlain isCompact>
                        <CardTitle style={{ textAlign: "center" }}>
                            <Icon size="xl" iconSize="lg">
                                <CodeBranchIcon color='blue' />
                            </Icon>
                            <h1><b>Resolution Time Bug SLO</b></h1>
                        </CardTitle>
                        <CardBody>
                            <Bullseye>
                                <Text style={{ textAlign: "center" }} component={TextVariants.p}>
                                    Aims to ensure that Blocker, Critical, and Major bugs are resolved in a reasonable period
                                </Text>
                            </Bullseye>
                        </CardBody>
                    </Card>
                </FlexItem>
                <Divider
                    orientation={{
                        default: 'vertical',
                    }}
                />
                <FlexItem flex={{ default: 'flex_1' }}>
                    <Card isPlain isCompact>
                        <CardTitle style={{ textAlign: "center" }}>
                            <Icon size="xl" iconSize="lg">
                                <UserIcon color='blue' />
                            </Icon>
                            <h1><b>Response Time Bug SLO</b></h1>
                        </CardTitle>
                        <CardBody>
                            <Bullseye>
                                <Text style={{ textAlign: "center" }} component={TextVariants.p}>
                                    Aims to ensure that Blocker and Critical bugs are assigned in the early phase of the bug's life
                                </Text>
                            </Bullseye>
                        </CardBody>
                    </Card>
                </FlexItem>
                <Divider
                    orientation={{
                        default: 'vertical',
                    }}
                />
                <FlexItem flex={{ default: 'flex_1' }}>
                    <Card isPlain isCompact>
                        <CardTitle style={{ textAlign: "center" }}>
                            <Icon size="xl" iconSize="lg">
                                <FilterIcon color='blue' />
                            </Icon>
                            <h1><b>Triage Time Bug SLO</b></h1>
                        </CardTitle>
                        <CardBody>
                            <Bullseye>
                                <Text style={{ textAlign: "center" }} component={TextVariants.p}>
                                    Aims to ensure that untriaged bugs are prioritized in the early phase of the bug's life
                                </Text>
                            </Bullseye>
                        </CardBody>
                    </Card>
                </FlexItem>
            </Flex>
        </CardBody>
    </Card>
);