namespace DF {
  namespace Entity {
    public abstract class Entity : OneOf.OneOfBase<SimpleEntity> {
      public static Entity Import(DF.Game.API.Data.Entity pb) {
        switch(pb.Type) {
          case DF.Game.API.Constants.EntityType.Tank:
            return new SimpleEntity(pb.EntityId);
          default:
            break;
        }
        throw new System.ArgumentException(
          System.String.Format("Input EntityType {0} is not recognized", pb.Type));
      }
    }

    // TODO(minkezhang): Rename to TANK eventually.
    public class SimpleEntity : Entity {
      private string _id;

      public SimpleEntity(string id) {
        _id = id;
      }

      public string ID { get => _id; }
    }
  }
}
