// <auto-generated>
//     Generated by the protocol buffer compiler.  DO NOT EDIT!
//     source: engine/server/executor/api/constants.proto
// </auto-generated>
#pragma warning disable 1591, 0612, 3021
#region Designer generated code

using pb = global::Google.Protobuf;
using pbc = global::Google.Protobuf.Collections;
using pbr = global::Google.Protobuf.Reflection;
using scg = global::System.Collections.Generic;
namespace DF.Game.Server.Service.API.Constants {

  /// <summary>Holder for reflection information generated from engine/server/executor/api/constants.proto</summary>
  public static partial class ConstantsReflection {

    #region Descriptor
    /// <summary>File descriptor for engine/server/executor/api/constants.proto</summary>
    public static pbr::FileDescriptor Descriptor {
      get { return descriptor; }
    }
    private static pbr::FileDescriptor descriptor;

    static ConstantsReflection() {
      byte[] descriptorData = global::System.Convert.FromBase64String(
          string.Concat(
            "CiplbmdpbmUvc2VydmVyL2V4ZWN1dG9yL2FwaS9jb25zdGFudHMucHJvdG8S",
            "IWdhbWUuc2VydmVyLnNlcnZpY2UuYXBpLmNvbnN0YW50cyp+CgxTZXJ2ZXJT",
            "dGF0dXMSGQoVU0VSVkVSX1NUQVRVU19VTktOT1dOEAASHQoZU0VSVkVSX1NU",
            "QVRVU19OT1RfU1RBUlRFRBABEhkKFVNFUlZFUl9TVEFUVVNfUlVOTklORxAC",
            "EhkKFVNFUlZFUl9TVEFUVVNfU1RPUFBFRBADKo4BCgxDbGllbnRTdGF0dXMS",
            "GQoVQ0xJRU5UX1NUQVRVU19VTktOT1dOEAASFQoRQ0xJRU5UX1NUQVRVU19O",
            "RVcQARIaChZDTElFTlRfU1RBVFVTX0RFU1lOQ0VEEAISFAoQQ0xJRU5UX1NU",
            "QVRVU19PSxADEhoKFkNMSUVOVF9TVEFUVVNfVEVBUkRPV04QBEJKWiFnYW1l",
            "LnNlcnZlci5zZXJ2aWNlLmFwaS5jb25zdGFudHOqAiRERi5HYW1lLlNlcnZl",
            "ci5TZXJ2aWNlLkFQSS5Db25zdGFudHNiBnByb3RvMw=="));
      descriptor = pbr::FileDescriptor.FromGeneratedCode(descriptorData,
          new pbr::FileDescriptor[] { },
          new pbr::GeneratedClrTypeInfo(new[] {typeof(global::DF.Game.Server.Service.API.Constants.ServerStatus), typeof(global::DF.Game.Server.Service.API.Constants.ClientStatus), }, null, null));
    }
    #endregion

  }
  #region Enums
  /// <summary>
  /// ServerStatus represents if the Executor currently is running the main tick
  /// loop. This is not surfaced in the client / server API.
  /// </summary>
  public enum ServerStatus {
    [pbr::OriginalName("SERVER_STATUS_UNKNOWN")] Unknown = 0,
    [pbr::OriginalName("SERVER_STATUS_NOT_STARTED")] NotStarted = 1,
    [pbr::OriginalName("SERVER_STATUS_RUNNING")] Running = 2,
    [pbr::OriginalName("SERVER_STATUS_STOPPED")] Stopped = 3,
  }

  /// <summary>
  /// ClientStatus represents the internal game's awareness of the current
  /// networking state of a connected client.
  /// </summary>
  public enum ClientStatus {
    [pbr::OriginalName("CLIENT_STATUS_UNKNOWN")] Unknown = 0,
    [pbr::OriginalName("CLIENT_STATUS_NEW")] New = 1,
    [pbr::OriginalName("CLIENT_STATUS_DESYNCED")] Desynced = 2,
    [pbr::OriginalName("CLIENT_STATUS_OK")] Ok = 3,
    [pbr::OriginalName("CLIENT_STATUS_TEARDOWN")] Teardown = 4,
  }

  #endregion

}

#endregion Designer generated code