export interface TestCase {
    name: string;
    test_case_impact: number;
    count: number;
    messages: {
        job_id: string;
        job_url: string;
        error_message: string;
        failure_date: string;
    }[]
}

export interface Flakey {
    status: string;
    test_cases: TestCase[];
    suite_name: string;
    average_impact: number;
}

export interface FlakeyObject {
    global_impact: number;
    git_organization: string;
    repository_name: string;
    job_name: string;
    suites: Flakey[];
}