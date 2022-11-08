import React, { useState, useLayoutEffect, useRef } from 'react';
import { ExclamationCircleIcon, OkIcon, HelpIcon } from '@patternfly/react-icons';
import { Card, CardTitle, CardBody } from '@patternfly/react-core';
import { Chart, ChartAxis, ChartBar, ChartLine, ChartGroup, ChartVoronoiContainer, ChartLegend, ChartLabel } from '@patternfly/react-charts';
import { SimpleList, SimpleListItem } from '@patternfly/react-core';
import { JobsComponent } from '@app/Jobs/Jobs';

/* 
Some common useful types definition
*/

export interface JobsStatistics {
  repository_name: string;
  git_organization: string;
  type: string;
  jobs: (JobsEntity)[] ;
}

export interface JobsEntity {
  html_url: string;
  name: string;
  metrics: (MetricsEntity)[];
  summary: (MetricsSummary);
}

export interface MetricsSummary {
  success_rate_avg: string;
  failure_rate_avg: string;
  ci_failed_rate_avg: string;
  date_from: string;
  date_to: string;
}

export interface MetricsEntity {
  success_rate: string;
  failure_rate: string;
  ci_failed_rate: string;
  date: string;
}

/* 
Simple list in a card, with selectable element
Accepts data in <SimpleListData> format type; requires a function to be executed onSelection
*/

export type SimpleListProps = {
  data: SimpleListData[],
  onSelection: (value) => void,
  title: any
};

export type SimpleListData = {
  index: number;
  value: string;
};

export const DashboardSimpleList = ({data, onSelection, title}:SimpleListProps) => {
  const onSelect = (selectedItem, selectedItemProps) => {
    onSelection(selectedItemProps["data-index"])
  }

  const items = data.map((job) => <SimpleListItem className="" key={job.index} data-index={job.index} isActive={job.index==0}> {job.value} </SimpleListItem>);

  return (
      <Card style={{width: "100%", height: "100%", fontSize: "1rem"}}>
        <CardTitle>{title}</CardTitle>
        <CardBody>
          <SimpleList onSelect={onSelect} aria-label="Simple List Example">
            {items}
          </SimpleList>
        </CardBody>
      </Card>
  )
};

/* 
Dashboard line chart
*/

export type DashboardLineChartDataSerie = {
  x: string;
  y: number;
  name: string;
}

export type DashboardLineChartSerie = {
  data: DashboardLineChartDataSerie[];
  style?: any;
}

export type DashboardLineChartData = {
  [key: string]: DashboardLineChartSerie
}

export const DashboardLineChart = ({data, colorScale}: {data:DashboardLineChartData, colorScale?:string[]}) => {
  const ref = useRef<HTMLDivElement>(null);
  let beautifiedData = data
  const [width, setWidth] = useState(0);
  const [height, setHeight] = useState(0);

  useLayoutEffect(() => {
    if(ref.current){
      setWidth(ref.current.offsetWidth-0);
      setHeight(ref.current.offsetHeight-40);
    }
  }, []);

  const legendData = Object.keys(data).map((key, index) => {
    let l = { name: key.toLowerCase().replace(/_/gi, " ") }
    if(data[key].style.data.stroke){
      l["symbol"] = {fill: data[key].style.data.stroke} 
    }
    
    return l
  })

  // Custom legend component
  const getLegend = (legendData) => (
    <ChartLegend
      borderPadding={{top: 20}}
      data={legendData}
      gutter={15}
      itemsPerRow={3}
      rowGutter={5}
    />
  );

  return (
    <div style={{ height: '100%', width: '100%', minHeight: "600px"}} className={"pf-c-card"} ref={ref}>
      <div style={{ height: height+'px', width: width+'px', background: "white" }}>
        <Chart
          ariaDesc="Average number of pets"
          ariaTitle="Bar chart example"
          containerComponent={<ChartVoronoiContainer labels={({ datum }) => `${datum.name}: ${datum.y}`} constrainToVisibleArea />}
          legendOrientation="horizontal"
          legendPosition="bottom"
          legendComponent={getLegend(legendData)}
          legendAllowWrap={true}
          height={height}
          padding={{
            bottom: 110,
            left: 80,
            right: 80,
            top: 80
          }}
          width={width}
        >
          <ChartAxis showGrid></ChartAxis>
          <ChartAxis dependentAxis showGrid />
          <ChartGroup>
            {
              Object.keys(data).map((key, index) => {
                return <ChartLine data={data[key].data} style={data[key].style} key={index}/>
              })
            }
          </ChartGroup>
        </Chart>
      </div>
    </div>
  )
};

/* 
Simple info card 
*/

export type InfoCardProp = {
  title: string;
  value: string;
}

export const InfoCard = ({data}: {data:InfoCardProp[]}) => {
  return (
    <Card style={{width: "100%", height: "100%", fontSize: "1rem"}}>
      <CardBody>
        { 
          data.map(function(value, index){
            return ( <div style={{marginTop: "5px"}} key={index}>
              <div><strong>{value.title}</strong></div>
              <div>{value.value}</div>
              </div> 
            )
          })
        }
      </CardBody>
    </Card>
  )
};

/* 
Dashboard Card with different style; will display a title and a value 
*/

export type DashboardCardProps = {
  cardType?: 'default' | 'danger' | 'success' | 'warning' | 'primary' | 'help';
  title: string;
  body: string;
};

export const DashboardCard = ({cardType, title, body}:DashboardCardProps) => {
  const cardStyle = new Map();
  cardStyle.set('title-danger', {color: "#A30000", fontWeight: "semibold", fontSize: "0.8em"});
  cardStyle.set('title-success', {color: "#1E4F18", fontWeight: "semibold", fontSize: "0.8em"});
  cardStyle.set('title-warning', {color: "#F0AB00", fontWeight: "semibold", fontSize: "0.8em"});
  cardStyle.set('title-default', {color: "black", fontWeight: "semibold", fontSize: "0.8em"});
  cardStyle.set('title-help', {color: "grey", fontWeight: "semibold", fontSize: "1em"});
  cardStyle.set('title-primary', {color: "#0066CC", fontWeight: "semibold", fontSize: "0.8em"});
  cardStyle.set('body-danger', {color: "#A30000", fontWeight: "bold", fontSize: "2em", textAlign: "center"});
  cardStyle.set('body-success', {color: "#1E4F18", fontWeight: "bold", fontSize: "2em", textAlign: "center"});
  cardStyle.set('body-warning', {color: "#F0AB00", fontWeight: "bold", fontSize: "2em", textAlign: "center"});
  cardStyle.set('body-default', {color: "black", fontWeight: "bold", fontSize: "1.8em", textAlign: "center"});
  cardStyle.set('body-help', {color: "grey", fontWeight: "normal", fontSize: "0.8em", textAlign: "left"});
  cardStyle.set('body-primary', {color: "#0066CC", fontWeight: "bold", fontSize: "2em", textAlign: "center"});

  return (
    <Card style={{width: "100%", height: "100%"}}>
      <CardTitle style={cardStyle.get("title-"+cardType)}>
        {cardType == 'help' && <HelpIcon style={{marginRight: "5px", fontSize: "1.1em", fontWeight: "bold", verticalAlign: "middle"}}></HelpIcon>}  
        {title}
      </CardTitle>
      <CardBody style={cardStyle.get("body-"+cardType)}>
        {body}
        {cardType == 'danger' && <ExclamationCircleIcon style={{fontSize: "1.2rem", margin: "0 5px"}}></ExclamationCircleIcon>}
        {cardType == 'success' && <OkIcon style={{fontSize: "1.2rem", margin: "0 5px"}}></OkIcon>}
        </CardBody>
    </Card>
  )
};