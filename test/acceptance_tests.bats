bats_require_minimum_version 1.5.0

setup() {
    load 'test_helper/bats-support/load'
    load 'test_helper/bats-assert/load'
    DIR="$( cd "$( dirname "$BATS_TEST_FILENAME" )" >/dev/null 2>&1 && pwd )"
    PATH="$DIR/..:$PATH"
}

@test 'display version' {
    run -- getignore --version
    assert_output --partial 'getignore version'
}

@test 'list files' {
    run getignore list
    assert_line 'C.gitignore'
    assert_line 'Global/Vim.gitignore'
    assert_line 'Yeoman.gitignore'
}

@test 'get file contents' {
    run getignore get C Global/Vim Yeoman
    assert_line '# C #'
    assert_line '# Vim #'
    assert_line '# Yeoman #'
}
