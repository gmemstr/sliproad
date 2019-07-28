import React, { Component } from 'react';
import { BrowserRouter, Route, Link } from 'react-router-dom'
import './App.css';

class App extends Component {
  render() {
    return (
      <BrowserRouter>
        <div className="App">
          <Route exact path="/" component={Homepage} />
          <Route exact path="/login" component={Login} />
          <Route path="/hot/:dir?" component={HotFileListing} />
        </div>
      </BrowserRouter>
    );
  }
}

class Login extends Component {
  render() {
    return (
      <div className="LoginForm">
        <form method="POST">
          <label>Username <input type="text" name="username"></input></label>
          <label>Password <input type="password" name="password"></input></label>
          <input type="submit" value="Login"></input>
        </form>
      </div>
    )
  }
}

class Homepage extends Component {
  constructor(props) {
    super(props);

    this.state = {
      diskusage: {},
      loading: true,
    };
  };

  componentDidMount() {
    fetch("/api/diskusage")
      .then(response => response.json())
      .then(data => this.setState({ diskusage: {
        hot: 100 - Math.floor((data.HotStorage.Free) / (data.HotStorage.Total) * 100),
        cold: 100 - Math.floor((data.ColdStorage.Free) / (data.ColdStorage.Total) * 100),
      }, loading: false }));
  }

  render() {
    const { loading } = this.state;

    if (loading) {
      return <div className="LoadingSpinner"><div></div><div></div><div></div><div></div></div>;
    }
    return (
      <div>
        <div className="Usages">
        <span>Hot Storage Usage: {this.state.diskusage.hot}%</span>
        <span>Cold Storage Usage: {this.state.diskusage.cold}%</span>
        </div>

        <div className="Navigation">
            <Link to="/hot">Hot</Link>
            <Link to="/cold">Cold</Link>
          </div>
      </div>
    )
  }
}

class HotFileListing extends Component {
  constructor(props) {
    super(props);

    this.state = {
      files: [],
      loading: true,
      directory: "",
    };
  };

  componentDidMount() {
    const { match: { params } } = this.props;
    this.setState({directory: params.dir});

    this.loadFileListing(params.dir);
  }

  componentDidUpdate(prevProps, prevState) {
    const { match: { params } } = this.props;
    console.log("Prev state: ", prevState);

    if (prevState.directory != params.dir) {
      this.setState({directory: params.dir});
      this.loadFileListing(params.dir);
    }
  }

  loadFileListing(dir) {
    if (dir != undefined) {
    fetch("/api/hot/" + dir)
      .then(response => response.json())
      .then(data => this.setState({ files: data, loading: false }));
    }
    else {
      fetch("/api/hot/")
      .then(response => response.json())
      .then(data => this.setState({ files: data, loading: false }));
    }
  }

  render() {
    const { loading } = this.state;

    if (loading) {
      return <div className="LoadingSpinner"><div></div><div></div><div></div><div></div></div>;
    }
    if (!loading && !this.state.files.Files) {
      return (
        <div>
          <FileUploadForm tier="hot" />
          Empty
        </div>
      )
    }

    return (
      <div>
          <FileUploadForm tier="hot" />
          <FileList files={this.state.files.Files} tier="hot" />
      </div>
    )
  }
}

class FileUploadForm extends Component {
  render() {
    return (
      <form className="FileUpload" enctype="multipart/form-data" method="POST" action={`/api/upload/${this.props.tier}`}>
        <input type="file" name="file" id="file" />
        <input type="submit" value="Upload" />
      </form>
    )
  }
}

class FileList extends Component {
  render () {
    console.log(this.props.files);
    let fileComponents = this.props.files.map((file) => {
      console.log(file)
      if (file.IsDirectory) {
        return <Directory dir={file} tier={this.props.tier} />
      }
      return <File file={file} tier={this.props.tier} />
    })
    return (
      <div>
        <ul>{fileComponents}</ul>
      </div>
    )
  }
}

class Directory extends Component {
  render() {
    return (
      <div>
        <Link to={`/${this.props.tier}/${this.props.dir.Name}`}>{this.props.dir.Name}/</Link>
      </div>
    )
  }
}

class File extends Component {
  render() {
    return (
      <div>
        <a href={`/api/${this.props.tier}/file/${this.props.file.Name}`}>{this.props.file.Name}</a>
      </div>
    )
  }
}

export default App;
