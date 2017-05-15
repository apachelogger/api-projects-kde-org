task :doc do
  system('node_modules/.bin/apidoc',
         '-e', 'node_modules',
         '-e', 'vendor',
         '-e', 'doc',
         '-o', 'doc', '--debug') || raise
end
