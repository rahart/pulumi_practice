from pulumi_policy import (
    EnforcementLevel,
    PolicyPack,
    ReportViolation,
    ResourceValidationArgs,
    ResourceValidationPolicy,
)
def resource_groups_tags_validator(args: ResourceValidationArgs, report_violation: ReportViolation):
    if args.resource_type == "pulumi:pulumi:Stack" or args.resource_type == "pulumi:providers:azure-native":
        # These dont seem to be azure resources so they are not tagable that.
        return

    req_tags = ["Environment", "Name", "Stack", "Project"]
    if "tags" in args.props:
        tags = args.props["tags"]
        for tag in req_tags:
            if tag not in tags:
                report_violation("The '{}' tag key must be added to the {}".format(tag, args.resource_type))
    else:
        report_violation(
            "Tags are required on all resources. " +
            "Please add the minimum of the following keys, " +
            "Required Keys: {}".format(', '.join(req_tags))
            )
def storage_account_no_public_read_validator(args: ResourceValidationArgs, report_violation: ReportViolation):
    if args.resource_type == "azure-native:storage:StorageAccount" and "allowBlobPublicAccess" in args.props:
        access_type = args.props["allowBlobPublicAccess"]
        if access_type == True: # TH: I set this a string at first. Very difficult/silly to debug.
            report_violation(
                "Azure Storage Account must not be set allowBlobPublicAccess: True."
            )

storage_acc_no_public_read = ResourceValidationPolicy(
    name="storage-account-no-public-read",
    description="Prohibits setting the public permission on Azure Storage Account.",
    validate=storage_account_no_public_read_validator,
)

require_resource_group_tag = ResourceValidationPolicy(
    name="require-resource-group-tagging",
    description="Requires the use of tags.",
    validate=resource_groups_tags_validator,
)

PolicyPack(
    name="azure-python",
    enforcement_level=EnforcementLevel.MANDATORY,
    policies=[
        storage_acc_no_public_read,
        require_resource_group_tag
    ],
)
