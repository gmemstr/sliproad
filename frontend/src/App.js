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
          <Route exact path="/hot" component={HotFileListing} />
        </div>
      </BrowserRouter>
    );
  }
}

class Login extends Component {
  render() {
    return (
      <div className="LoginForm">
        <form>
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
    };
  };

  componentDidMount() {
    fetch("/api/hot/")
      .then(response => response.json())
      .then(data => this.setState({ files: data, loading: false }));
  }

  render() {
    const { loading } = this.state;

    if (loading) {
      return <div className="LoadingSpinner"><div></div><div></div><div></div><div></div></div>;
    }

    return (
      <div>
          <FileList files={this.state.files.Files} />
      </div>
    )
  }
}

class FileList extends Component {
  render () {
    let fileComponents = this.props.files.map((file) => {
      if (file.IsDirectory) {
      return <Directory dir={file}/>
      }
      return <File file={file}/>
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
        <p>{this.props.dir.Name}/</p>
      </div>
    )
  }
}

class File extends Component {
  render() {
    return (
      <div>
        <p>{this.props.file.Name}</p>
      </div>
    )
  }
}

export default App;
