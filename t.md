# [CORPORATE / BUSINESS] PROCEDURE

## Document Information

| Field | Value |
|-------|-------|
| **Document Library:** | Platform Engineering / Serverless Platform Team |
| **Page:** | 1 of 4 |
| **Title:** | Application Container Service (ACS) Customer Onboarding Procedure |
| **Version:** | [1.0] Policy Office to populate |
| | **Scheduled Review:** [Annually] |
| | **Owned by:** | Platform Engineering Team |
| | **Parent Policy/Standards:** | [IT Infrastructure Management Policy, API Management Standards, Containerized Application Deployment Standards] |
| | **Date Approved:** MM/DD/YYYY \| **Effective Date:** MM/DD/YYYY |

---

## Scope

This Procedure applies to GEICO Corporation (the "Company"), its direct and indirect Controlled Subsidiaries (together with the Company, the "GEICO Group"), and all employees of any of the Companies. The Company reserves the right to change, modify, add, or remove portions of this Procedure at any time.

---

## Purpose

This procedure enables internal GEICO development teams to successfully onboard to the Application Container Service (ACS) platform for deploying containerized applications. Teams will learn how to configure deployment specifications, integrate with the orchestration API, establish automated CI/CD pipelines, and validate successful deployments on the ACS platform.

---

## Procedure Details

This procedure describes the complete process for onboarding development teams to the Application Container Service (ACS) platform, from initial prerequisites through successful deployment validation. The process ensures teams can deploy modern, scalable, containerized applications while maintaining GEICO's security and operational standards.

### I. Roles and Responsibilities

**A. Platform Engineering Team**

   i. **Provide technical guidance** and onboarding support for customer teams
      • Maintain ACS platform documentation and templates
      • Respond to technical questions via #help-paas-serverless Slack channel

   ii. **Manage platform infrastructure** including namespace provisioning and cluster management

   iii. **Maintain CI/CD templates** for both Azure DevOps and GitHub platforms

**B. Customer Development Teams**

   i. **Complete prerequisite requirements** before beginning ACS onboarding

   ii. **Configure deployment specifications** according to platform standards

   iii. **Set up automated CI/CD pipelines** using provided templates

**C. Supporting Teams**

   i. **PATH Team:** Provides System Entry IDs for application registration

   ii. **Cluster Management Team:** Manages Kubernetes infrastructure and namespace provisioning

   iii. **GEICO Secrets Platform (GSP) Team:** Provides secret management capabilities

### II. Prerequisites and Requirements

**A. System Registration (PATH)**

   i. **Obtain System Entry ID** from the System Catalog (PATH) portal
      • Navigate to PATH Developer Portal Systems section
      • Create new system entry or locate existing system
      • Record System Entry ID for use in deployment configurations

   ii. **Verify system status** shows as "In Service" in PATH portal

**B. Repository Setup**

   i. **Establish code repository** using either GitHub or Azure DevOps
      • Ensure contributor access to chosen repository
      • Clone repository locally and verify code deployment readiness

   ii. **Verify repository contains** containerized application ready for deployment

**C. Platform Identity Creation**

   i. **Configure authentication credentials** for CI/CD pipeline execution
      • Follow PaaS Identity Creation process for required access permissions
      • Locate ADO Service Connection or GitHub Environment Variable configuration

   ii. **Verify identity has necessary permissions** to interact with ACS orchestration API

**D. Secret Management (Optional)**

   i. **If application requires secret management,** complete GSP onboarding process
      • Contact Secret Management Team at SecretsMngmntTeam@geico.com
      • Gather required information: Cluster Name, Namespace Name, Service Account Name, Team AD Group
      • Follow GSP User Guide for onboarding instructions

### III. Deployment Configuration

**A. Create Configuration Directory Structure**

   i. **Create config folder** in application root directory using `mkdir config`

   ii. **Create environment-specific spec file** named `spec.<env>.yaml` where `<env>` represents target environment

**B. Select Appropriate Spec File Template**

   i. **For new customers (onboarded after 05/19/2025):** Use NaaS-integrated template
      • Includes automated namespace provisioning
      • Supports networkZone parameter (st=Semi-Trusted, ut=Untrusted)
      • Provides enhanced automation capabilities

   ii. **For existing customers (legacy):** Use backward-compatible template
      • Requires pre-provisioned namespaces
      • Contains manual cluster and ingress IP configuration
      • Note: Legacy templates will be deprecated

**C. Configure Spec File Parameters**

   i. **Update required metadata fields:**
      • systemEntryId: Use System ID from PATH
      • componentName: Must start with lowercase letter, alphanumeric only
      • notifications: Include team email addresses

   ii. **Configure placement settings:**
      • For new customers: Set networkZone (st or ut)
      • For legacy customers: Specify cluster, configKey, and ingressIP

   iii. **Define container specifications:**
      • Set resource requests and limits (CPU, memory)
      • Configure health check probes (liveness, readiness)
      • Specify scaling parameters (min/max replicas, metrics)

### IV. CI/CD Pipeline Setup

**A. Choose Pipeline Platform**

   i. **Azure DevOps (ADO) Setup:**
      • Create `.cicd` directory in project root
      • Create `publish-release.yml` file using provided ADO template
      • Configure pool, variables, and job parameters according to template

   ii. **GitHub Actions Setup:**
      • Create `.github` directory in project root
      • Create `workflows` subdirectory
      • Create `publish-release.yml` file using provided GitHub template

**B. Configure Pipeline Variables**

   i. **Update template-specific variables:**
      • Image name and tag configuration
      • JFrog repository paths and credentials
      • System Entry ID and component name
      • Environment and region specifications

   ii. **Set deployment parameters:**
      • Environment selection (sb, dv, it, in, ut, pp, pd, tr, np, ft)
      • Region specification (east, west, plaza, fredericksburg)
      • Action type (deploy, delete, switchover, podrestart, poddelete)

### V. Deployment Execution and Validation

**A. Execute Initial Deployment**

   i. **For ADO pipelines:** Navigate to Azure Pipelines and select "New Pipeline"

   ii. **For GitHub pipelines:** Commit workflow file to trigger automated execution

   iii. **Monitor pipeline execution** for successful completion and artifact publishing

**B. Verify Deployment in ArgoCD**

   i. **Access ArgoCD dashboard** using AzureAD authentication
      • OnPrem ArgoCD Deployments for on-premises clusters
      • ACS ArgoCD Deployments for Azure clusters

   ii. **Search for deployment** using PATH System ID

   iii. **Validate deployment status** shows as "Healthy" and "Synced"

   iv. **Verify essential information:**
      • Application health and sync status
      • Repository and target revision information
      • Namespace and cluster destination

### VI. Application Testing and Validation

**A. Connectivity Testing**

   i. **Identify external dependencies** required by the application

   ii. **Test connectivity** using tools like ping, telnet, or curl

   iii. **Submit firewall access requests** if connectivity tests fail

**B. Functional Testing**

   i. **Web Applications:**
      • Retrieve endpoint URL from deployment pipeline output
      • Perform validation using browser or testing tools (curl, Insomnia)
      • Verify all functionalities work as expected

   ii. **Worker Roles:**
      • Monitor deployment and application logs for proper execution
      • Manually trigger jobs or tasks to verify worker functionality
      • Validate outputs and results meet expectations

   iii. **CronJobs:**
      • Verify CronJob scheduling configuration and execution intervals
      • Monitor execution logs to confirm scheduled operation
      • Validate task results and outputs

### VII. Troubleshooting and Support

**A. Common Issues Resolution**

   i. **Deployment failures:** Check ArgoCD sync status and event logs

   ii. **Connectivity issues:** Verify network security group configurations

   iii. **Authentication problems:** Validate platform identity and service connections

**B. Support Channels**

   i. **Primary support:** #help-paas-serverless Slack channel

   ii. **Documentation:** ACS User Guide FAQ section

   iii. **Escalation:** Platform Engineering Team for complex technical issues

---

## Defined Terms

**Application Container Service (ACS):** GEICO's next-generation serverless platform for deploying and managing containerized applications with automated scaling, security, and high availability.

**ArgoCD:** A declarative, GitOps continuous delivery tool for Kubernetes that automates application deployment and manages DNS creation.

**GEICO Secrets Platform (GSP):** GEICO's centralized platform for managing application secrets and encryption services.

**NaaS (Namespace as a Service):** Automated namespace provisioning service that dynamically creates and manages Kubernetes namespaces based on application requirements.

**PATH:** GEICO's System Catalog platform used for application registration and System Entry ID management.

**Spec File:** YAML configuration template containing application host configurations, runtime settings, and service integration specifications required for ACS deployment.

**System Entry ID:** Unique identifier assigned by PATH for application registration and deployment tracking purposes.

---

## Enforcement

All employees whose responsibilities are affected by this Procedure are expected to be familiar with the details and responsibilities created by this document. Failure to comply with this Procedure may result in disciplinary action, up to and including termination of employment.

To report an issue of non-compliance or procedure violation, contact the Policy Owner, Process Owner, Policy Office, or submit an **Ethics Violation Form**.

---

## References & Exhibits

### [Related Documents]

[ACS User Guide - Deployment Configuration](internal-link)

[PaaS Identity Creation Guide](internal-link)

[GSP User Guide](internal-link)

[System Catalog (PATH) Documentation](internal-link)

### [Exhibit 1: Deployment Flow Diagram]

*[A visual representation of the complete onboarding process from prerequisites through validation would be included here]*

### [Exhibit 2: Template Selection Matrix]

| Customer Type | Onboarding Date | Template Type | NaaS Integration | Status |
|---------------|----------------|---------------|------------------|---------|
| New Customers | After 05/19/2025 | Enhanced | Yes | Active |
| Legacy Customers | Before 05/19/2025 | Traditional | No | Deprecated |

---

## Revision History

| Field | Value |
|-------|-------|
| **Revisions** | MM/DD/YYYY (1.0) - Initial procedure creation |
| | **(** **[ ]** = Version # |
| | *(# after decimal increased +1 since revision was made within the same year)* |
| | *(# before decimal increased +1 since revision was made within a separate year)* |
| **Next Scheduled Review** | MM/DD/YYYY |
| | *[This line will be updated to account for a rolling 12-month review cycle. It will be 12 months from last Approved Date]* |

*\*Format only, no substantive changes made.*
