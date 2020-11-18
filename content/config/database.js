module.exports = ({ env }) => ({
  defaultConnection: 'default',
  connections: {
    default: {
      connector: 'bookshelf',
      settings: {
        client: 'postgres',
        host: env('DATABASE_HOST', '127.0.0.1'),
        port: env.int('DATABASE_PORT', 5433),
        database: env('DATABASE_NAME', 'strapi_dev'),
        username: env('DATABASE_USERNAME', 'strapi'),
        password: env('DATABASE_PASSWORD', 'magical_password'),
        ssl: env.bool('DATABASE_SSL', false),
      },
      options: {}
    },
  },
});
