var gulp         = require('gulp');
var pogo         = require('gulp-pogo');
var browserify   = require('browserify');
var sass         = require('gulp-sass');
var concat       = require('gulp-concat');
var plumber      = require('gulp-plumber');
var gutil        = require('gulp-util');
var fs           = require('fs');
var watch        = require('gulp-watch');
var karma        = require('karma').server;
var watchify     = require('watchify');
var source       = require('vinyl-source-stream');
var markdownify  = require('markdownify');
var argv         = require('yargs').argv;
var uglify       = require('gulp-uglify');
var gulpif       = require('gulp-if');
var buffer       = require('vinyl-buffer');
var size         = require('gulp-size');
var autoprefixer = require('gulp-autoprefixer');

var onError = function (err) {
  gutil.beep();
  gutil.log(gutil.colors.red(err.message))
  gutil.log(err)
};

gulp.task('styles', function (callback) {
  return gulp.src('./public/css/app.scss')
    .pipe(plumber({
      errorHandler: onError
    }))
    .pipe(sass())
    .pipe(autoprefixer())
    .pipe(concat('bundle.css'))
    .pipe(gulp.dest('./public'))
});

/**
 * Run test once and exit
 */
gulp.task('test', function (done) {
  karma.start({
    configFile: __dirname + '/karma.conf.js',
    singleRun: true
  }, done);
});

/**
 * Watch for file changes and re-run tests on each change
 */
gulp.task('tdd', function (done) {
  karma.start({
    configFile: __dirname + '/karma.conf.js'
  }, done);
});

gulp.task("watchify", function() {
    browserifyAndMaybeWatchify(true)
})

gulp.task("browserify", ["compile-pogo"], function() {
    browserifyAndMaybeWatchify(false)
})

gulp.task('compile-pogo', function(callback){
    return gulp.src('./public/{js,test}/**/*.pogo')
      .pipe(plumber({errorHandler: onError}))
      .pipe(pogo())
      .pipe(gulp.dest('./public/'));
})

function browserifyAndMaybeWatchify(watch) {
  args = watchify.args;
  args.extensions = ['.md'];

  var bundler = browserify("./public/js/app.js", args);

  bundler.transform(markdownify);

  var bundle = function() {
    return bundler
      .bundle()
      .on('error', onError)
      .pipe(source('bundle.js'))
      .pipe(buffer())
      .pipe(gulpif(argv.production, uglify()))
      .pipe(size())
      .pipe(gulp.dest('./public/'));
  };

  if (watch) {
    bundler = watchify(bundler);
    bundler.on("update", bundle);
  }

  bundle()
}

gulp.task("watch", ["compile-pogo", "watchify"], function() {
  watch('./public/{js,test}/**/*.pogo')
    .pipe(plumber({errorHandler: onError}))
    .pipe(pogo())
    .pipe(gulp.dest('./public/'));

  gulp.watch('./public/css/**', ['styles']);
})

gulp.task('default', ['styles', 'browserify']);
