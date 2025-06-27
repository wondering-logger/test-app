Week 1: Foundation Tests (Start Here)
Container Apps Basic CRUD:
gherkin✅ Test_ContainerApp_CreateAndDeploy_BasicWeb
   → Create simple web app → Verify namespace created → Verify pods running → Verify HTTP endpoint

✅ Test_ContainerApp_CreateAndDeploy_BasicWorker
   → Create worker app → Verify deployment → Verify background processing

✅ Test_ContainerApp_CreateAndDeploy_BasicCronJob
   → Create cronjob → Verify schedule → Verify execution → (You already have this!)

🆕 Test_ContainerApp_Delete_CleanupResources
   → Create app → Delete app → Verify namespace deleted → Verify no orphaned resources

🆕 Test_ContainerApp_Update_ResourceQuotas
   → Create app → Update CPU/memory → Verify namespace quota updated
NaaS Integration Basics:
gherkin🆕 Test_NaaS_NamespaceCreation_NewCustomer
   → Deploy with networkZone → Verify namespace auto-created → Verify naming convention

🆕 Test_NaaS_NamespaceCreation_LegacyCustomer
   → Deploy with pre-existing namespace → Verify no NaaS call → Verify deployment succeeds
Week 2-3: Error Scenarios & Edge Cases
Authorization & Validation:
gherkin🆕 Test_ContainerApp_Create_InvalidSystemEntryId
   → Try create with bad systemEntryId → Verify 400 error → Verify no resources created

🆕 Test_ContainerApp_Delete_UnauthorizedUser
   → Create app as user A → Try delete as user B → Verify 403 → Verify app still exists

🆕 Test_ContainerApp_Create_DuplicateComponentName
   → Create app → Try create same name different owner → Verify 400 error
Service Integration Failures:
gherkin🆕 Test_ContainerApp_Deploy_NaaSServiceDown
   → Mock NaaS unavailable → Try deploy → Verify graceful failure → Verify rollback

🆕 Test_ContainerApp_Deploy_JFrogAuthFailure
   → Invalid JFrog credentials → Try deploy → Verify clear error message

🆕 Test_ContainerApp_Deploy_InvalidSpecFile
   → Upload malformed spec → Try deploy → Verify validation errors
Namespace Edge Cases:
gherkin🆕 Test_ContainerApp_Delete_NamespaceStuckTerminating
   → Create app → Force namespace into terminating state → Try delete → Verify error handling

🆕 Test_ContainerApp_Deploy_ResourceQuotaExceeded
   → Request huge resources → Try deploy → Verify quota error → Verify no partial deployment
Week 3-4: Advanced Scenarios
Multi-Region & Cloud Providers:
gherkin🆕 Test_ContainerApp_Deploy_MultiRegion
   → Create app in east+west → Verify both deployments → Verify separate namespaces

🆕 Test_ContainerApp_Deploy_AzureVsOnPrem
   → Create identical apps on Azure and on-prem → Verify different cluster targeting

🆕 Test_ContainerApp_Delete_MultiRegion
   → Create multi-region app → Delete → Verify all namespaces cleaned up
Deployment Strategies:
gherkin🆕 Test_ContainerApp_Deploy_BlueGreenStrategy
   → Deploy v1 → Deploy v2 with blue-green → Verify traffic switching → Verify old version cleanup

🆕 Test_ContainerApp_Deploy_CanaryStrategy
   → Deploy v1 → Deploy v2 with canary → Verify traffic split → Verify gradual rollout
Database & State Management:
gherkin🆕 Test_ContainerApp_Deploy_DatabaseConsistency
   → Create app → Verify database records match cluster state → Update app → Verify consistency

🆕 Test_ContainerApp_Restart_PreserveState
   → Create app → Restart pods → Verify database state unchanged → Verify app recovers
Month 2: Performance & Reliability
Load & Performance:
gherkin🆕 Test_ContainerApp_Deploy_HighConcurrency
   → Deploy 50 apps simultaneously → Verify all succeed → Verify no resource conflicts

🆕 Test_ContainerApp_Deploy_LargeApplication
   → Deploy app with many containers → Verify reasonable deployment time → Verify resource usage

🆕 Test_ContainerApp_Scaling_AutoScale
   → Deploy app → Generate load → Verify auto-scaling → Verify scale-down
Disaster Recovery:
gherkin🆕 Test_ContainerApp_Deploy_ClusterFailover
   → Deploy to primary cluster → Simulate cluster failure → Verify failover to secondary

🆕 Test_ContainerApp_Deploy_ServiceMeshResilience
   → Deploy with service mesh → Simulate network partitions → Verify recovery
Integration Complexity:
gherkin🆕 Test_ContainerApp_Deploy_FullStack
   → Deploy app with secrets, volumes, service mesh, monitoring → Verify all integrations

🆕 Test_ContainerApp_Deploy_SecretsManagement
   → Deploy with GSP secrets → Verify secret injection → Update secrets → Verify app updates
