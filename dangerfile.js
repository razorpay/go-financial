import { danger, fail, warn,  message } from "danger";
const pr = danger.github.pr;
const modifiedFiles = danger.git.modified_files.concat(danger.git.created_files);

/**
 * Rule: Small pr is suggested.
 * Reason: Pr is not supposed to be very large so it is suggested to keep the pr small.
 *         So it would be easy for reviewer to review and reduce chances of missing bugs.
 *         Warn when there is a big PR
 */
(async function () {
    await checkPRSize();
})();

async function checkPRSize() {
    const bigPRThresholdWarn = 300;
    const bigPRThresholdFail = 400;
    const ignoredLineCount = await getIgnoredLineCount();

    if (pr.additions - ignoredLineCount > bigPRThresholdFail && pr.labels.filter(label => label.name  == "skip_pr_size_check").length == 0 ) {
        fail(
            `Your PR has over ${bigPRThresholdFail} lines of code additions :scream: . Try to breakup into separate PRs :+1:`
        );
    }else if (pr.additions - ignoredLineCount> bigPRThresholdWarn) {
        warn(
            `Your PR has over ${bigPRThresholdWarn} lines of code additions :scream: . Try to breakup into separate PRs :+1:`
        );
    }

}


async function getIgnoredLineCount() {
    let ignoredLineCount = 0;
    const testChanges = modifiedFiles.filter(filepath =>
        filepath.includes('test'),
    );

    await Promise.all(
        testChanges.map(async (file) => {
            const diff = await danger.git.structuredDiffForFile(file);
            diff.chunks.map((chunk) => {
                // Here we filter to only get the additions
                const additions = chunk.changes.filter(({ type }) => type === "add");
                ignoredLineCount = ignoredLineCount + additions.length;
            });
        })
    );

    return ignoredLineCount;
}

/**
 * Rule: Pr description is required.
 * Reason: No PR is too small to include a description of why you made a change
 *         1931 is size of pr template assuming that is never going to change
 */
const prTemplateSize = 50;
if (pr.body == null || pr.body.length < prTemplateSize) {
    warn(`Please include a description of your PR changes.`);
}



/**
 * Rule: Any change to source file require changes to test file also.
 * Reason: No code changes are good if you don't add proper test cases. This might result in
 *         big issues in production and require extra work of debugging
 */

const internalAppChanges = modifiedFiles.filter(filepath =>
    filepath.includes('internal'),
);
const pkgAppChanges = modifiedFiles.filter(filepath =>
    filepath.includes('pkg'),
);
const testChanges = modifiedFiles.filter(filepath =>
    filepath.includes('test'),
);

const appChanges = internalAppChanges.length + pkgAppChanges.length;
if (appChanges > 0 && testChanges.length < 1) {
    warn ("Remember to write tests in case you have added a new API or fixed a bug. Feel free to ask for help if you need it 👍");
} else if (testChanges.length > 0) {
    message(`Thanks 🙏 for adding test cases you are awesome 😎`);
}


/**
 * Rule: Small pr is suggested by number of file changed.
 * Reason: Pr is not supposed to be very large so it is suggested to keep the pr small.
 *         So it would be easy for reviewer to review and reduce chances of missing bugs.
 *         Warn when there is a big PR
 */

const excludedFilesChanges = modifiedFiles.filter(filepath =>
    filepath.includes('go.mod') ||
    filepath.includes('go.sum') ||
    filepath.includes('config/') ||
    filepath.includes('test')
);

const fileChangedThreshold = 15;
if (pr.changed_files - excludedFilesChanges.length > fileChangedThreshold) {
    warn(
        `Your PR has over ${fileChangedThreshold} changes file :scream: . Try to breakup into separate PRs :+1:`
    );
}

